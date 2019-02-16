/* 2019-02-16 (cc) <paul4hough@gmail.com>
   FIXME what is this for?
*/
package ticket

func (t *Ticket) manager(tsys uint64) {

	for {
		for _, tq := range t.db.TicketQueue(tsys) {
			t.ProcQueue(tsys,tq)
		}
		sleep(10 minutes)
	}
}

func (t *Ticket)ProcQueue(tsys, tq uint64) {

	for _, ta := range t.db.TicketActions(tsys,tq) {
		if ! t.ProcAction(tsys,tq,ta) {
			return
		}
	}
}

func (t *Ticket)ProcAction(tsys, tq, ta uint64) bool {

	tact := t.db.TicketActionGet(tsys, tq, ta)
	if tact == nil {
		return true
	}

	switch tact.action {
	case Create:
		tckt := t.db.TicketGetTicket(tsys, tq)
		tid := t.ActionCreate(tsys,tckt.group,tckt.title,tact.desc)
		if len(tid) > 0 {
			if t.db.TicketSetTid(tsys,tq,tid) {
				return t.db.TicketActionDel(tsys, tq, ta)
			} else {
				return false
			}
		} else {
			return false
		}
	case Update:
		tid := t.db.TicketGetTid(tsys, tq)
		if  t.ActionUpdate(tsys,tid,tact.payload) {
			return t.db.TicketActionDel(tsys, tq, ta)
		} else {
			return false
		}
	case Close:
		tid := t.db.TicketGetTid(tsys, tq)
		if ! t.ActionClose(tsys,tid,tact.payload) {
			return t.db.TicketDelete(tsys, tq, ta)
		} else {
			return false
		}
	}
}
