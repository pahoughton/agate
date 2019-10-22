/* 2019-10-21 (cc) <paul4hough@gmail.com>
   notify note
*/
package note

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	pmod "github.com/prometheus/common/model"
)

const (
	LABS_TITLE = "title"
	TIMEFMT = "2006-01-02 15:04:05.9999 -0700"
)

type Alert struct {
	Name	string
	Labels	pmod.LabelSet
	Annots	pmod.LabelSet
	Starts	time.Time
	From	string
	Labsfp	pmod.Fingerprint
}

func (self *Alert) Text(prefix string) string {
	text := prefix + "from: " + self.From + "\n"
	text += prefix + self.Name + " " + self.Starts.Format(TIMEFMT) + "\n\n"
	for k, v := range self.Labels {
		text += fmt.Sprintf("%s%16s: %s\n",prefix,k,v)
	}
	// need?
	if len(self.Annots) > 0 {
		text += prefix + "annotations\n"
		for k, v := range self.Annots {
			text += fmt.Sprintf("%s%16s: %s\n",prefix,k,v)
		}
	}
	return text
}

type Note struct {
	Labels  pmod.LabelSet
	Alerts  []Alert
	From	string
	Nid		[]byte
}

func (self *Note) Title() string {

	var title string
	if t, ok := self.Labels[LABS_TITLE]; ok {
		title = string(t)
	} else {
		for _, lk := range SortedKeys(self.Labels) {
			title += string(self.Labels[pmod.LabelName(lk)]) + " "
			if len(title) > 70 {
				break
			}
		}
	}

	if len(self.Alerts) > 1 {
		title +=  " " + strconv.Itoa(len(self.Alerts)) + " alerts"
	}
	if len(title) > 0 {
		return title
	} else {
		return "Unknown alert"
	}
}

func (self *Note) Desc() string {

	text := "from: " + self.From + "\n"
	for _, lk := range SortedKeys(self.Labels) {
		text += fmt.Sprintf("  %16s: %s\n",lk,self.Labels[pmod.LabelName(lk)])
	}
	if len(self.Labels) > 0 {
		text += "\n"
	}
	for _, a := range self.Alerts {
		text += a.Text("") + "\n"
	}
	return text
}

func (self *Note) String() string {
	return self.Title() + "\n" + self.Desc()
}

// return len == 0 no changes
func (self *Note) Changes(n []Alert) string {

	nMap := make(map[pmod.Fingerprint]Alert,len(n))
	oMap := make(map[pmod.Fingerprint]Alert,len(self.Alerts))

	var text string
	for _, v := range n {
		nMap[v.Labsfp] = v
	}

	tmp := ""
	for _, v := range self.Alerts {
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


// todo general library function
func SortedKeys(m pmod.LabelSet) []string {

	keys := make([]string, 0, len(m))
	for k, _ := range m {
		keys = append(keys,string(k))
	}
	sort.Strings(keys)
	return keys
}
