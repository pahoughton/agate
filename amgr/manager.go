/* 2019-02-14 (cc) <paul4hough@gmail.com>

Single AlertGroup Queue Manager Thread
*/
package amgr

func (am *Amgr)Manager() {

	for {
		// grab array of queue keys
		agq := a.db.AGroupQueue()
		if len(agq) < 1 {
			// wait for next alert, double check queue every 10 min
			select {
			case <- h.manager:
			case <- time.After(10 * time.Minute):
			}
			agq = a.db.AGroupQueue()
		}

		for _, agqkey := range agq {
			am.procq <- agkey
			go ProcAlertGroup(agqkey)
		}
	}
}

type NewAlert struct {
	agidx	int
	remed	bool
}

func (am *Amgr)ProcAlertGroup(agqkey uint64) {

	defer func() {
		<- am.procq
	}
	agq := a.db.AGroupGet(agqkey)
	if agq == nil {
		return
	}
	if a.debug {
		var dbgbuf bytes.Buffer
		if err := json.Indent(&dbgbuf, agq.Json, "", "  "); err != nil {
			fmt.Printf("DEBUG json.Indent: ",err.Error())
		} else {
			fmt.Println("DEBUG agrp\n",dbgbuf.String())
		}
	}

	var agrp alert.AlertGroup
	if err := json.Unmarshal(agq.Json, &agrp); err != nil {
		panic(fmt.Sprintf(
			"json.Unmarshal agrp: %s\n%v",err.Error(),agrcv.Json))
	}
	if len(agrp.Alerts) < 1 {
		am.Error("0 alerts in alertgroup")
		a.db.AGroupDelete(agqkey)
		return
	}
	var gtid *db.AlertTicket

	newAlerts := make([]NewAlert,0,len(agrp.Alerts))
	resolvedAlerts := make([]int,0,len(agrp.Alerts))
	anyRemed := false

	for agidx, a := range agrp.Alerts {

		atid := am.db.AlertTicketGet(a.StartsAt, a.Key())
		if atid != nil {
			gtid = atid
		}
		if a.Status == 'firing' {
			if atid != nil {
				continue
			}
			aname := a.Name()
			resolve := "false"
			if ag.Resolve {
				resolve = "true"
			}
			am.metrics.AlertsFiring.With(
				promp.Labels{
					"name": aname,
					"node": a.Node(),
					"resolve": resolve
				}).Inc()

			remed := false
			if aname == "unknown" {
				am.Error("alert missing alertname")
			} else {
				sfn := path.Join(h.proc.ScriptsDir,aname)
				finfo, err = os.Stat(sfn)
				if err == nil && (finfo.Mode() & 0111) != 0 {
					remed = true
				} else {
					ardir := path.Join(h.proc.PlaybookDir,"roles",aname)
					finfo, err := os.Stat(ardir)
					if err == nil && finfo.IsDir() {
						remed = true
					}
				}
			}
			anyRemed ||= remed
			newAlerts = append(newAlerts,
				&NewAlert{
					agidx: agidx,
					remed: remed,
				})
		} else if a.Status == 'resolved' {
			if atid != nil {
				resolvedAlerts = append(resolvedAlerts,agidx)
			}
		} else {
			am.Error("unknown status: " + a.Status " "+ a.Title())
		}
	}

	if gtid == nil {
		gtid = am.ticket.AGroupCreate(agrp,agq.resolv,anyRemed)
	} else {
		for _, a := range newAlerts {
			am.ticket.AGroupAppendAlert(gtid,agrp[a.agidx],a.remed)
		}
	}
	for _, a := range newAlerts {
		if a.remed {
			am.ticket.AGroupAppend(
				gtid,
				agrp[a.agidx],
				am.Remediate(agrp[a.agidx]))
		}
	}
	for _, agidx := range resolvedAlerts {
		am.ticket.AGroupResolved(agrp[agidx])
	}
	a.db.AGroupDelete(agqkey)
}
