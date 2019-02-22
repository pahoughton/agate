/* 2019-02-19 (cc) <paul4hough@gmail.com>

respond to an amgr alerg group

- retrieve json
- unmarshal
- alerts:
  - ? firing:
    - ? new
    - ? Labels[ticket_sys]
    - ? Labels[ticket_grp]
    - ? Remediate
  - ? resolved
    - ? new
- ? new alert group
  - ticket.Create
  - db.update
- ? new alerts
  - ticket.Update
  - db.update
- ? remediate
  - fix
- ? new resolved
  - ticket.Update
  - db.update
- ? all resolved
  - ticket.Close
*/
package amgr

type NewAlert struct {
	agidx	int
	remed	bool
}

func (am *Amgr)Respond(agqkey uint64) {

	// retrieve json
	ag := am.db.AGroupGet(agqkey)
	if ag == nil {
		return
	}
	if a.debug {
		var dbgbuf bytes.Buffer
		if err := json.Indent(&dbgbuf, ag.Json, "", "  "); err != nil {
			fmt.Printf("DEBUG json.Indent: ",err.Error())
		} else {
			fmt.Println("DEBUG agrp\n",dbgbuf.String())
		}
	}
	// unmarshal
	var agrp alert.AlertGroup
	if err := json.Unmarshal(ag.Json, &agrp); err != nil {
		panic(fmt.Sprintf(
			"json.Unmarshal agrp: %s\n%v",err.Error(),agrcv.Json))
	}
	if len(agrp.Alerts) < 1 {
		panic("0 alerts in alertgroup")
	}

	tid := (*tid.Tid)nil
	newAlerts := make([]NewAlert,0,len(agrp.Alerts))
	resolvedAlerts := make([]int,0,len(agrp.Alerts))
	anyRemed := false
	resolveCount := 0

	tsyscnt = make(map[string]int,len(agrp.Alerts))
	tgrpcnt = make(map[string]int,len(agrp.Alerts))

	// alerts:
	for agidx, a := range agrp.Alerts {

		atid := am.db.AlertTid(a.StartsAt, a.Key())
		// ? new alert group
		if atid != nil {
			tid = atid
		}
		// ? firing:
		if a.Status == 'firing' {

			aname := a.Name()

			// ? new
			if atid != nil {
				continue
			} else if tid == nil {
				// ? Labels[ticket_sys]
				if v, ok := a.Labels["ticket_sys"]; ok {
					tsyscnt[string(v)] += 1
				}
				// ? Labels[ticket_grp]
				if v, ok := a.Labels["ticket_grp"]; ok {
					tgrpcnt[string(v)] += 1
				}
			}
			resolve := "false"
			if agq.Resolv {
				resolve = "true"
			}
			am.metrics.Alerts.With(
				promp.Labels{
					"name": aname,
					"node": a.Node(),
					"resolve": resolve
				}).Inc()

			// ? Remediate
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
			// ? resolved
			resolveCount += 1
			if atid != nil {
				// ? new
				resolvedAlerts = append(resolvedAlerts,agidx)
			}
		} else {
			am.Error("unknown status: " + a.Status " "+ a.Title())
		}
	}

	// ? ticket.Create
	if tid == nil {

		// ticket system(gitlab,mock,...) to use
		tsys := am.t.Default
		if v, ok := agrp.ComLabels["ticket_sys"]; ok {
			tsys = ticket.NewTSys(string(v))
		} else {
			majTSys := am.t.Default.String()
			majCount := 0
			for k, c := range tsyscnt {
				if c > majCount {
					majTSys = ticket.NewTSys(k)
					majCount = c
				}
			}
			tsys = majTSys
		}
		// ticket group to use
		tgrp := am.t.Group(tsys)
		if v, ok :=  agrp.ComLabels["ticket_grp"]; ok {
			tgrp = v
		} else {
			majTGrp := tgrp
			majCount := 0
			for k, c := range tgrpcnt {
				if c > majCount {
					majTGrp = k
					majCount = c
				}
			}
			tgrp = majTGrp
		}
		hdr := "\n"
		if anyRemed {
			hdr += "remediation: pending\n\n"
		} else {
			hdr += "remediation: none\n\n"
		}
		if agq.Resolve {
			hdr += "close: auto\n\n"
		} else {
			hdr += "close: manual\n\n"
		}
		tid = am.t.Create(tsys, tgrp, ag.Title(), hdr + ag.Desc())
		// db.update
		if ag.Resolve {
			for _, a := range agrp.Alerts {
				am.db.AlertAdd(a.StartsAt,a.Key(),tid.Bytes())
			}
		}
	} else if len(newAlerts) > 0 {
		// - ? new alerts
		msg := "\nNew Alerts\n"
		for _, a := range newAlerts {
			msg += agrp.Alerts[a.agidx].Title() + "\n"
			msg += agrp.Alerts[a.agidx].Desc() + "\n"
		}
		// ticket.Update
		am.ticket.Update(tid, msg)
	}
	// ? remediate
	for _, a := range newAlerts {
		if a.remed {
			// fix
			am.ticket.Update(tid,am.remed.Fix(agrp[a.agidx]))
		}
	}
	// new resolved
	if len(resolvedAlerts) > 0 {
		msg := "\n"
		for _, agidx := range resolvedAlerts {
			a = agrp.Alerts[agidx]
			msg += "resolved: " + a.Title() + "\n"
			// db.update
			am.db.AlertDel(a.StartsAt,a.Key())
		}
		// ticket.Update
		am.ticket.Update(tid,msg)
		// ? all resolved
		if resolveCount >= len(agrp.Alerts) {
			am.ticket.Update(tid,"\nAll Alerts Resolved\n")
			if am.ticket.CloseResolved {
				// ticket.Close
				am.ticket.Close(tid)
		}
	}
	// - db update
	a.db.AGroupDel(agqkey)
}
