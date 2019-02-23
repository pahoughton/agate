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

import (
	// "bytes"
	"encoding/json"
	"fmt"
	promp "github.com/prometheus/client_golang/prometheus"
	"github.com/pahoughton/agate/amgr/alert"
	// "github.com/pahoughton/agate/ticket"
	"github.com/pahoughton/agate/ticket/tid"
)

type NewAlert struct {
	agidx	int
	remed	bool
}

func (am *Amgr)Respond(agqkey uint64) bool {

	// retrieve json
	ag := am.db.AGroupGet(agqkey)
	if ag == nil {
		return true
	}
	/*
	if am.debug {
		var dbgbuf bytes.Buffer
		if err := json.Indent(&dbgbuf, ag.Json, "", "  "); err != nil {
			fmt.Println("DEBUG json.Indent: ",err.Error())
		} else {
			fmt.Println("DEBUG agrp\n",dbgbuf.String())
		}
	}
    */
	// unmarshal
	var agrp alert.AlertGroup
	if err := json.Unmarshal(ag.Json, &agrp); err != nil {
		panic(fmt.Sprintf(
			"json.Unmarshal agrp: %s\n%v",err.Error(),ag.Json))
	}
	if len(agrp.Alerts) < 1 {
		panic("0 alerts in alertgroup")
	}

	var agtid tid.Tid
	newAlerts := make([]NewAlert,0,len(agrp.Alerts))
	resolvedAlerts := make([]int,0,len(agrp.Alerts))
	anyRemed := false
	resolveCount := 0

	tsyscnt := make(map[string]int,len(agrp.Alerts))
	tgrpcnt := make(map[string]int,len(agrp.Alerts))

	// alerts:
	for agidx, a := range agrp.Alerts {
		//fmt.Printf(" alert %v %v\n",a.StartsAt,a.Key())
		atid := am.db.AlertGet(a.StartsAt, a.Key())
		// ? new alert group
		if atid != nil {
			agtid = tid.NewBytes(atid)
		}
		// ? firing:
		if a.Status == "firing" {

			aname := a.Name()

			// ? new
			if atid != nil {
				continue
			} else if agtid == nil {
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
			if ag.Resolve {
				resolve = "true"
			}
			am.metrics.alerts.With(
				promp.Labels{
					"name": aname,
					"node": a.Node(),
					"resolve": resolve,
				}).Inc()

			// ? Remediate
			remed := false
			if aname == "unknown" {
				am.Errorf("alert missing alertname")
			} else {
				remed = remed || am.remed.AnsibleAvail(a.Labels)
				remed = remed || am.remed.ScriptAvail(a.Labels)
			}
			anyRemed = anyRemed || remed
			newAlerts = append(newAlerts,
				NewAlert{
					agidx: agidx,
					remed: remed,
				})

		} else if a.Status == "resolved" {
			// ? resolved
			resolveCount += 1
			if atid != nil {
				// ? new
				resolvedAlerts = append(resolvedAlerts,agidx)
			}
		} else {
			am.Errorf("unknown status: %v - %v",a.Status,a.Title())
		}
	}
	// ? ticket.Create
	if agtid == nil {
		// ticket system(gitlab,mock,...) to use
		tsysStr := am.ticket.Default.String()
		if v, ok := agrp.ComLabels["ticket_sys"]; ok {
			tsysStr = string(v)
		} else {
			majTSys := tsysStr
			majCount := 0
			for k, c := range tsyscnt {
				if c > majCount {
					majTSys = k
					majCount = c
				}
			}
			tsysStr = majTSys
		}
		tsys := am.ticket.NewTSysString(tsysStr)
		// ticket group to use
		tgrp := am.ticket.Group(tsys)
		if v, ok :=  agrp.ComLabels["ticket_grp"]; ok {
			tgrp = string(v)
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
		if ag.Resolve {
			hdr += "close: auto\n\n"
		} else {
			hdr += "close: manual\n\n"
		}
		agtid = am.ticket.TCreate(tsys,tgrp,agrp.Title(),hdr + agrp.Desc())
		if agtid == nil {
			return false// abort
		}
		// db.update - need for dup detection
		for _, a := range agrp.Alerts {
			am.db.AlertAdd(a.StartsAt,a.Key(),agtid.Bytes())
		}
	} else if len(newAlerts) > 0 {
		// - ? new alerts
		msg := "\nNew Alerts\n"
		for _, a := range newAlerts {
			msg += agrp.Alerts[a.agidx].Title() + "\n"
			msg += agrp.Alerts[a.agidx].Desc() + "\n"
		}
		// ticket.Update
		if am.ticket.TUpdate(agtid, msg) == false {
			return false
		}
	}
	// ? remediate
	for _, a := range newAlerts {
		if a.remed {
			// fix
			out := am.Fix(agrp.Alerts[a.agidx])
			if am.ticket.TUpdate(agtid,out) == false {
				am.Errorf("remed ticket(%s) update\n%v",agtid.String(),out)
				return false
			}
		}
	}
	// new resolved
	if len(resolvedAlerts) > 0 {
		msg := "\n"
		for _, agidx := range resolvedAlerts {
			a := agrp.Alerts[agidx]
			msg += "resolved: " + a.Title() + "\n"
			// db.update
			am.db.AlertDel(a.StartsAt,a.Key())
		}
		// ? all resolved
		if resolveCount >= len(agrp.Alerts) {
			msg += "\nAll Alerts Resolved\n"
			if am.ticket.CloseResolved {
				// ticket.Close
				if am.ticket.TClose(agtid,msg) == false {
					return false
				}
			} else {
				if am.ticket.TUpdate(agtid,msg) == false {
					return false
				}
			}
		} else {
			// ticket.Update
			if am.ticket.TUpdate(agtid,msg) == false {
				return false
			}
		}
	}
	// - db update
	am.db.AGroupDel(agqkey)
	return true
}
