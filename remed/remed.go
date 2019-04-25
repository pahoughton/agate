/* 2019-02-14 (cc) <paul4hough@gmail.com>
   run alert remediation
*/
package remed

import (
	"sync/atomic"
	"github.com/pahoughton/agate/amgr/alert"
	"github.com/pahoughton/agate/notify/nid"

)


func (r *Remed) remed(a alert.Alert,nid nid.Nid) {
	defer r.wg.Done()
	defer atomic.AddInt32(&r.cnt,-1)

	defer r.metrics.remedq.Dec()
	r.metrics.remedq.Inc()

	out := ""

	if r.AnsibleAvail(a.LabelSet()) {
		out += "ansible remed:"
		tmp, err := r.Ansible(a.Node(),a.LabelSet())
		if err != nil {
			r.error(err)
			out += " ERROR - " + err.Error()
		}
		if len(tmp) > 0 {
			out += "\n" + tmp
		}
	}
	if r.ScriptAvail(a.LabelSet()) {
		out += "script remed:"
		tmp, err := r.Script(a.Node(),a.LabelSet())
		if err != nil {
			r.error(err)
			out += " ERROR - " + err.Error()
		}
		if len(tmp) > 0 {
			out += "\n" + tmp
		}
	}
	if len(out) < 1 {
		out = a.Name() + " no remed output"
		r.errorf(out)
	}
	if r.notify.Update(nid,out) == false {
		r.errorf("remed notify(%s) update\n%v",nid.Id(),out)
	}
}

func (r *Remed) AlertHasRemed(a alert.Alert) bool {
	return r.ScriptAvail(a.LabelSet()) || r.AnsibleAvail(a.LabelSet())
}

func (r *Remed) AGroupHasRemed(ag alert.AlertGroup) bool {
	for _, a := range ag.Alerts {
		if r.AlertHasRemed(alert.Alert(a)) {
			return true
		}
	}
	return false
}
func (r *Remed) Remed(a alert.Alert, nid nid.Nid) {
	if r.AlertHasRemed(a) {
		atomic.AddInt32(&r.cnt,1)
		if r.cnt >= r.parallel {
			r.wg.Wait()
		}
		r.wg.Add(1)
		go r.remed(a,nid)
	}
}
