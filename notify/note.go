/* 2019-10-19 (cc) <paul4hough@gmail.com>
   FIXME what is this for?
*/
package notify


type Note struct {
}

func (n *Note) Title() string {

}

func (n *Note) Desc() string {

}

func (o *Note) Changes(n []Alerts) string {

	nMap := make(map[pmod.Fingerprint]bool,len(n))
	oMap := make(map[pmod.Fingerprint]bool,len(o.Alerts))

	text := "\nupdates:\n"
	for _, v := range n {
		nMap[v.labsfp] = true
	}

	rList := make([]Alerts,len(a))
	tmp := ""
	for _, v := range o.Alerts {
		if _, ok := nMap[v.labsfp]; ! ok {
			tmp += a.Text("    ") + "\n"
		}
		oMap[v.labsfp] = true
	}
	if len(tmp) > 0 {
		text += "\n  resolved:\n" + tmp + "\n"
	}

	for k, v := range nMap {
		if _, ok := oMap[v.labsfp]; ! ok {
			tmp += a.Text("    ") + "\n"
		}
	}
	if len(tmp) > 0 {
		text += "\n  firing:\n" + tmp + "\n"
	}
	return text
}
