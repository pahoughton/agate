/* 2019-02-14 (cc) <paul4hough@gmail.com>
   run alert remediation
*/
package remed

import (
	"sync/atomic"
	"strings"
	pmod "github.com/prometheus/common/model"
	"github.com/pahoughton/agate/notify"

)
var (
	NODE_LABELS = []string{"agate_node", "hostname", "node", "instance"}
)

func labelSetNode(l pmod.LabelSet) string {
	keys := NODE_LABELS
	for _, k := range keys {
		if v, ok := l[pmod.LabelName(k)]; ok {
			node := string(v)
			if i := strings.IndexRune(node,':'); i > 0 {
				return node[:i]
			} else {
				return node
			}
		}
	}
	return ""
}

//func (r *Remed) remed(a alert.Alert,nid nid.Nid) {
func (r *Remed) remed(task string, labels pmod.LabelSet, nkey notify.Key) {

	defer r.wg.Done()
	defer atomic.AddInt32(&r.cnt,-1)

	defer r.metrics.remedq.Dec()
	r.metrics.remedq.Inc()

	out := ""

	if r.AnsibleAvail(task) && len(labelSetNode(labels)) > 0 {
		out += "ansible remed:"
		tmp, err := r.Ansible(task,labelSetNode(labels),labels)
		if err != nil {
			r.error(err)
			out += " ERROR - " + err.Error()
		}
		if len(tmp) > 0 {
			out += "\n" + tmp
		}
	}
	if r.ScriptAvail(task) {
		out += "script remed:"
		tmp, err := r.Script(task,labels)
		if err != nil {
			r.error(err)
			out += " ERROR - " + err.Error()
		}
		if len(tmp) > 0 {
			out += "\n" + tmp
		}
	}
	if len(out) < 1 {
		out = task + " no remed output"
		r.errorf(out)
	}
	if r.notify.UpdateNote(nkey,out) == false {
		r.errorf("remed notify(%s) update\n%v",nkey,out)
	}
}

func (r *Remed) HasRemed(task string) bool {
	return r.ScriptAvail(task) || r.AnsibleAvail(task)
}

func (r *Remed) Remed(task string, labels pmod.LabelSet, nkey notify.Key) {
	if r.HasRemed(task) {
		atomic.AddInt32(&r.cnt,1)
		if r.cnt >= r.parallel {
			r.wg.Wait()
		}
		r.wg.Add(1)
		go r.remed(task, labels, nkey)
	}
}
