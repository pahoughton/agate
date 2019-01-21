/* 2019-01-19 (cc) <paul.houghton.ywi9@statefarm.com>
   agate models
*/
package model

import (
	"sort"
	"strings"

	pmod "github.com/prometheus/common/model"
)

type Alert struct {
	pmod.Alert

	Status	pmod.AlertStatus	`json:"status"`
}

type AlertGroup struct {

	Version		string			`json:"version"`
	Receiver	string			`json:"receiver"`
	Status		string			`json:"status"`
	ExtURL		string			`json:"externalURL"`
	GroupKey	string			`json:"groupKey"`
	Alerts		[]Alert			`json:"alerts"`
	ComAnnots	pmod.LabelSet	`json:"commonAnnotations,omitempty"`
	ComLabels	pmod.LabelSet	`json:"commonLabels,omitempty"`
	GroupLabels	pmod.LabelSet	`json:"groupLabels,omitempty"`
}

func (a *Alert) Key() uint64 {
	return uint64(a.Labels.Fingerprint())
}

func (a *Alert) Title() string {

	var (
		title	string
		ok		bool
	)
	node := "unknown"

	if inst, ok := a.Labels["instance"]; ok {
		node = strings.Split(string(inst),":")[0]
	} else if h, ok := a.Labels["hostname"]; ok {
		node = string(h)
	}


	if _, ok = a.Annotations["title"]; ok {
		title = string(a.Annotations["title"])
	} else if  _, ok = a.Annotations["subject"]; ok {
		title = string(a.Annotations["subject"])
	} else {
		title = a.Name() + " on " + node
	}
	return title
}

func (a *Alert) Desc() string {

	desc := "from: " + a.GeneratorURL + "\n"
	desc += "when: " + a.StartsAt.String() + "\n"

	desc += "\nAnnotations:\n"
	ankeys := make([]string, 0, len(a.Annotations))
	for ak, _ := range a.Annotations {
		ankeys = append(ankeys, string(ak))
	}
	sort.Strings(ankeys)
	for _, ak := range ankeys {
		desc += ak + ": " +  string(a.Annotations[pmod.LabelName(ak)])  + "\n"
	}

	desc  += "\nLabels:\n"
	lbkeys := make([]string, 0, len(a.Labels))
	for lk, _ := range a.Labels {
		lbkeys = append(lbkeys, string(lk))
	}
	sort.Strings(lbkeys)
	for _, lk := range lbkeys {
		desc += lk + ": " + string(a.Labels[pmod.LabelName(lk)]) + "\n"
	}

	return desc
}
