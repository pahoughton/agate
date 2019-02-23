/* 2019-02-14 (cc) <paul4hough@gmail.com>
   run alert remediation
*/
package amgr

import (
	"github.com/pahoughton/agate/amgr/alert"
)
func (am *Amgr)Fix(a alert.Alert) string {

	out := ""

	if am.remed.AnsibleAvail(a.Labels) {
		if tmp, err := am.remed.Ansible(a.Node(),a.Labels); err != nil {
			am.Error(err)
		} else {
			out += tmp
		}
	}
	if am.remed.ScriptAvail(a.Labels) {
		if tmp, err := am.remed.Script(a.Node(),a.Labels); err != nil {
			am.Error(err)
		} else {
			out += tmp
		}
	}
	return out
}
