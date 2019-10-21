/* 2019-10-21 (cc) <paul4hough@gmail.com>
   notify note
*/
package note

import (
	"time"

	pmod "github.com/prometheus/common/model"
)

type Alert struct {
	Name	string
	Labels	pmod.LabelSet
	Annots	pmod.LabelSet
	Starts	time.Time
	Genurl	string
	Labsfp	pmod.Fingerprint
}

func (self *Alert) Text(prefix string) string {
	return "FIXME STUB"
}

type Note struct {
	Labels  pmod.LabelSet
	Alerts  []Alert
	From	string
	Updates string
	Nid		[]byte
}

func (n *Note) Title() string {
	return "FIXME STUB"
}

func (n *Note) Desc() string {
	return "FIXME STUB"
}

func (o *Note) Changes(n []Alert) string {

	nMap := make(map[pmod.Fingerprint]Alert,len(n))
	oMap := make(map[pmod.Fingerprint]Alert,len(o.Alerts))

	text := "\nupdates:\n"
	for _, v := range n {
		nMap[v.Labsfp] = v
	}

	tmp := ""
	for _, v := range o.Alerts {
		if _, ok := nMap[v.Labsfp]; ! ok {
			tmp += v.Text("    ") + "\n"
		}
		oMap[v.Labsfp] = v
	}
	if len(tmp) > 0 {
		text += "\n  resolved:\n" + tmp + "\n"
	}

	for k, v := range nMap {
		if _, ok := oMap[k]; ! ok {
			tmp += v.Text("    ") + "\n"
		}
	}
	if len(tmp) > 0 {
		text += "\n  firing:\n" + tmp + "\n"
	}
	return text
}
