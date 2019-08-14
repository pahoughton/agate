/* 2019-02-19 (cc) <paul4hough@gmail.com>
*/
package amgr

import (
	"fmt"
	"strconv"

	"github.com/pahoughton/agate/db"
	"github.com/pahoughton/agate/amgr/alert"
	"github.com/pahoughton/agate/notify"
	"github.com/pahoughton/agate/notify/nid"
)

func (am *Amgr)Notify(nsysid notify.NSys, qid uint64) bool {

	ag := alert.NewAlertGroup(am.db.AGroupQueueGet(uint(nsysid),qid))
	if ag == nil {
		panic("alertgroup nil")
	}
	nsys := am.db.AGroupNSysGet(ag.StartsAt(),ag.Key())
	if nsys == nil {
		if am.debug {
			fmt.Printf("debug race nsys nil for ag: %v\n",ag.Data)
		}
		// race - duplicate ag w/ deleted twin
		am.db.AGroupQueueDel(uint(nsysid),qid)
		return true
	}
	if am.respond(*nsys,*ag) {
		am.db.AGroupQueueDel(uint(nsysid),qid)
		return true
	} else {
		return false
	}
}

func (am *Amgr)respond(nsys db.NSys, ag alert.AlertGroup) bool {

	if bnid := am.db.AGroupNidGet(ag.StartsAt(),ag.Key()); bnid == nil {
		if ag.Status != "firing" {
			// new resolved alertgroup - ignore it
			return true
		}
		nid, err := am.notify.Create(
			notify.NSys(nsys.Sys),
			nsys.Grp,
			ag.Title(),
			ag.Desc(),
			am.remed.AGroupHasRemed(ag),
			nsys.Resolve)
		if err != nil {
			fmt.Printf("warn create fail: %s retry: %v err: %v\n",
				notify.NSys(nsys.Sys).String(),
				am.retry,
				err)
			return false
		}
		am.db.AGroupNidAdd(ag.StartsAt(),ag.Key(),nid.Bytes())

		// be sure to add alert records
		for _, a := range ag.Alerts {
			am.db.AlertNidAdd(a.StartsAt,alert.Alert(a).Key(),nid.Bytes())
		}
		// remed may block
		for _, a := range ag.Alerts {
			am.remed.Remed(alert.Alert(a),nid)
		}
	} else {
		nid := nid.NewBytes(bnid)
		update := ""
		if ag.Status == "firing" {
			for _, aga := range ag.Alerts {
				a := alert.Alert(aga)
				if am.db.AlertNidGet(a.StartsAt,a.Key()) == nil {
					update += "\nfiring: " + a.Title() + "\n" + a.Desc() + "\n"
					am.db.AlertNidAdd(a.StartsAt,a.Key(),nid.Bytes())
					am.remed.Remed(a,nid)
				}
			}
		} else if ag.Status == "resolved" {
			rcnt := 0
			ncnt := 0
			for _, aga := range ag.Alerts {
				a := alert.Alert(aga)
				if a.Status != "resolved" {
					continue
				} else {
					rcnt += 1
				}

				if v := am.db.AlertNidGet(a.StartsAt,a.Key()); v != nil {
					ncnt += 1
					update += "\nresolved: " + a.Title()
					am.db.AlertNidDel(a.StartsAt,a.Key())
				}
			}
			if rcnt == len(ag.Alerts) {
				if am.notify.CloseResolved {
					err := am.notify.Close(nid,"\nall resolved:\n" + update)
					if err != nil {
						fmt.Printf("warn close fail: %s retry: %v err: %v\n",
							notify.NSys(nsys.Sys).String(),
							am.retry,
							err)
						return false
					}
				}
				am.db.AGroupNSysDel(ag.StartsAt(),ag.Key())
				am.db.AGroupNidDel(ag.StartsAt(),ag.Key())
				return true
			} else {
				update = strconv.Itoa(ncnt) + " alerts resolved:\n" + update
			}
		} else {
			panic("unk status: " + ag.Status )
		}
		if len(update) > 0 {
			err := am.notify.Update(nid,update)
			if err != nil {
				fmt.Printf("warn update fail: %s retry: %v err: %v\n",
					notify.NSys(nsys.Sys).String(),
					am.retry,
					err)
				return false
			}
		}
	}
	return true
}
