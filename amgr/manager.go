/* 2019-02-14 (cc) <paul4hough@gmail.com>
   FIXME what is this for?
*/
package amgr

func (a *Amgr)Manager() {

	for {
		// wait for next alert, double check queue every 10 min
		select {
		case recvd := <- h.manager:
		case <- time.After(10 * time.Minute):
		}

		for {
			agrcv := a.db.AGroupNext()
			if id == nil {
				break
			}
			// fixme - one at a time - need rate limiter
			a.ProcAlertGroup(agrcv)
		}
	}
}

type RemedAlert struct {
	alert	*model.Alert
	remed	bool
}

func (am *Amgr)ProcAlertGroup(agrcv *db.AGroupRcvd) {

	if a.debug {
		var dbgbuf bytes.Buffer
		if err := json.Indent(&dbgbuf, agrcv.Json, "", "  "); err != nil {
			fmt.Printf("DEBUG json.Indent: ",err.Error())
		} else {
			fmt.Println("DEBUG agrp\n",dbgbuf.String())
		}
	}

	var agrp model.AlertGroup
	if err := json.Unmarshal(ag.Json, &agrp); err != nil {
		panic(fmt.Sprintf(
			"json.Unmarshal agrp: %s\n%v",err.Error(),agrcv.Json))
	}

	var gtid *db.AlertTicket

	if agrp.Status == 'firing' {

		ticketAlerts := make([]RemedAlert)
		remedAlerts := make([]Alert)

		for _, a := range agrp.Alerts {

			if a.Status == 'firing' {

				atid := am.db.TicketGet(a.StartsAt, a.Key())
				if len(atid) > 0 {
					gtid = atid
				} else {
					aname = a.Name()
					if agrcv.Resolve {
						am.metrics.AlertsFiring.With(
							promp.Labels{
								"name": aname,
								"node": a.Node(),
								"resolve": "true",
							}).Inc()
					} else {
						am.metrics.AlertsFiring.With(
							promp.Labels{
								"name": aname,
								"node": a.Node(),
								"resolve": "false",
							}).Inc()
					}
					remed := false
					if aname == "unknown" {
						am.Error("alert missing alertname")
					} else if a.Node() != "unknown" {
						sfn := path.Join(h.proc.ScriptsDir,aname)
						finfo, err = os.Stat(sfn)
						if err == nil && (finfo.Mode() & 0111) != 0 {
							remedAlerts = append(remedAlerts, a)
							remed = true
						} else {
							ardir := path.Join(
								h.proc.PlaybookDir,"roles",aname)
							finfo, err := os.Stat(ardir)
							if err == nil && finfo.IsDir() {
								remedAlerts = append(remedAlerts, a)
								remed = true
							}
						}
					}
					ra := &RemedAlert{
						alert: a,
						remed: remed,
					}
					ticketAlerts = append(ticketAlerts,ra)
				}
			}
		}
		if gtid != nil {
			for _, a := range ticketAlerts {
				am.ticket.Append(gtid,ra.alert,ra.remed)
				if agrcv.Resolve {
					am.db.TicketAdd(a.StartsAt,a.Key(),gtid)
				}
			}
		} else {
			gtid = am.ticket.Create(agrp,len(remedAlerts) > 0)
			if agrcv.Resolve {
				for _, a := range agrp.Alerts {
					am.db.TicketAdd(a.StartsAt,a.Key(),gtid)
				}
			}
		}
		for _, a := range remedAlerts {
			am.Remediate(a,gtid,len(agrp.Alerts) > 1)
		}
	} else {
		resolved := true

		for _, a := range agrp.Alerts {
			if a.Status == 'firing' {
				resolved = false
				continue
			} else {
				aKey := a.Key()
				gtid = am.db.TicketGet(a.StartsAt,aKey)
				if gtid != nil {
					am.metrics.AlertsResolved.With(
						promp.Labels{
							"name": a.Name(),
							"node": a.Node(),
						}).Inc()
					am.ticket.Resolved(a,tid)
					am.db.TicketDel(a.StartsAt,aKey)
				}
			}
		}
		if resolved && gtid != nil {
			am.ticket.Close(gtid)
		}
	}
}
