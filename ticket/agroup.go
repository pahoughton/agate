/* 2019-02-16 (cc) <paul4hough@gmail.com>
   Alert Group Ticket interface
*/
package ticket

func (t *Ticket)AGroupCreate(
	agrp alert.AlertGroup,resolve,remed bool) *db.AlertTid {

	gcnt := len(agrp.Alerts)
	mtsys := make([]string,len(agrp.Alerts))
	mtgrp := make([]string,len(agrp.Alerts))
	tsys := t.TsysDefault
	tgrp := t.TgrpDefault

	for agidx, a := range agrp.Alerts {
		if v, ok := a.Labels[TicketSys]; ok {
			mtsys = append(mtsys,v)
			tsys = v
		}
		if v, ok := a.Labels[TicketGrp]; ok {
			mtgrp = append(mtgrp,v)
			tgrp = v
		}
	}
	// multiple tsys or tgrp - report error and use first
	if len(mtsys) > 1 {
		t.Error(fmt.Sprintf("multiple tsys %v",mtsys))
		tsys := mtsys[0]
	}
	if let(mtgrp) > 1 {
		t.Error(fmt.Sprintf("multiple tgrp %v",mtgrp))
		tgrp := mtgrp[0]
	}

	sys := t.Sys(tsys)
	title := ""
	remedStr := ""
	desc := ""

	if remed {
		remedStr = "\nRemediation Pending\n"
	} else {
		remedStr = "\nNO remediation available\n"
	}
	if gcnt == 1 {
		title = agrp.Alerts[0].Title()
		desc = remedStr + agrp.Alerts[0].Desc()
	} else {
		title = agrp.Title()
		desc = remedStr + agrp.Desc()
	}

	tid,err := t.ActionCreate(sys,tgrp,title,desc)

	atid := &db.AlertTid{}

	if err != nil {
		atid.Tqid = = t.db.TicketQueueCreate(sys,tgrp,title,desc)
	}
	atgrp := make([]db.AlertKey,0,len(agrp.Alerts))
	for agidx, a := range agrp.Alerts {
		atgrp = append(atgrp,db.AlertKey{Start: a.StartsAt, Key: a.Key()})
	}
	t.db.AlertGroupTicketCreate(atgrp,atid)
	return atid
}

func (t *Ticket) AGroupAppendAlert(
	atid *db.AlertTid, a *alert.Alert,remed bool) {

	if cfg.NewTicketOnAdded {

	} else {

	}
}
