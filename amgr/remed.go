/* 2019-02-14 (cc) <paul4hough@gmail.com>
   run alert remediation
*/
package amgr

import (
	"sync/atomic"
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/pahoughton/agate/amgr/alert"
	"github.com/pahoughton/agate/ticket/tid"
)

func (am *Amgr) remed(a alert.Alert,tid tid.Tid) {
	defer am.fix.wg.Done()
	defer am.metrics.remed.With(prom.Labels{"alert": a.Name()}).Dec()

	am.metrics.remed.With(prom.Labels{"alert": a.Name()}).Inc()

	out := ""

	if am.fix.remed.AnsibleAvail(a.Labels) {
		if tmp, err := am.fix.remed.Ansible(a.Node(),a.Labels); err != nil {
			am.Error(err)
		} else {
			out += tmp
		}
	}
	if am.fix.remed.ScriptAvail(a.Labels) {
		if tmp, err := am.fix.remed.Script(a.Node(),a.Labels); err != nil {
			am.Error(err)
		} else {
			out += tmp
		}
	}
	if len(out) > 0 {
		if am.ticket.Update(tid,out) == false {
			am.Errorf("remed ticket(%s) update\n%v",tid.String(),out)
		}
	}
}

func (am *Amgr) Remed(a alert.Alert,tid tid.Tid) {
	atomic.AddInt32(&am.fix.cnt,1)

	if am.fix.cnt >= int32(am.fix.max) {
		am.fix.wg.Wait()
	}
	am.fix.wg.Add(1)
	go am.remed(a,tid)

}
