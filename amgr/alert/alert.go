/* 2019-01-19 (cc) <paul.houghton.ywi9@statefarm.com>
   agate models
*/
package alert

import (
	"encoding/binary"
	"fmt"
	"sort"
	"strings"
	// "time"
	pmod "github.com/prometheus/common/model"
)

type LabelMap pmod.LabelSet

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
	ComAnnots	LabelMap	`json:"commonAnnotations,omitempty"`
	ComLabels	LabelMap	`json:"commonLabels,omitempty"`
	GroupLabels	LabelMap	`json:"groupLabels,omitempty"`
}

var (
	ProcAnnotKeys = map[pmod.LabelName]bool {
		"agate_group_title":	true,
		"group_title":			true,
		"agate_title":			true,
		"metric":				true,
		"title":				true,
		"subject":				true,
	}
	ProcLabelKeys = map[pmod.LabelName]bool {
		"agate_node":	true,
		"hostname":		true,
		"instance":		true,
		"mongrp":		true,
	}
)
func (a *Alert) Key() []byte {
	if b, err := a.StartsAt.MarshalBinary(); err == nil {
		k := make([]byte,binary.MaxVarintLen64,
			binary.MaxVarintLen64 +
			len(b)+(len(b) % binary.MaxVarintLen64))
		binary.PutUvarint(k,uint64(a.Fingerprint()))
		return(append(k,b...))
	} else {
		panic(err)
	}

	/*
	k := make([]byte,0,binary.Size(time.Time{}) + binary.Size(a.Fingerprint))

		panic(err)
	} else {
		binary.PutUvarint(k,uint64(a.Fingerprint()))
		return append(k,b...)
	}
*/
}

func (lm *LabelMap)SortedKeys() pmod.LabelNames {

	lkeys := make(pmod.LabelNames, 0, len(*lm))
	for k, _ := range *lm {
		lkeys = append(lkeys, k)
	}

	sort.Sort(lkeys)
	return lkeys
}

func (a *Alert) Title() string {

	var (
		node	string
		title	string
	)

	keys := pmod.LabelNames{"agate_title", "title", "subject"}
	for _, k := range keys {
		if _, ok := a.Annotations[k]; ok {
			title = string(a.Annotations[k])
			break
		}
	}
	if len(title) == 0 {
		keys = pmod.LabelNames{"agate_node", "hostname"}
		for _, k := range keys {
			if _, ok := a.Labels[k]; ok {
				node = string(a.Labels[k])
				break
			}
		}
		if len(node) == 0 {
			if tmp, ok := a.Labels["instance"]; ok {
				node = strings.Split(string(tmp),":")[0]
			}
		}
		aname := a.Name()
		if len(aname) == 0 {
			aname = "unknown"
		}
		if len(node) > 0 {
			title = aname + " on " + node
		} else {
			title = aname
		}
	}
	return title
}

func (ag *AlertGroup) Title() string {
	var title string
	if tmp, ok := ag.ComAnnots["agate_group_title"]; ok {
		title = string(tmp)
	} else {
		title = fmt.Sprintf("%d grouped alerts",len(ag.Alerts))
	}
	return title
}

func (a *Alert) Desc() string {

	desc := "\nFrom: " + a.GeneratorURL + "\n\n"
	desc += "When: " + a.StartsAt.String() + "\n"

	keys := make(pmod.LabelNames, 0, len(a.Annotations))
	for k, _ := range a.Annotations {
		if ! ProcAnnotKeys[k] {
			keys = append(keys, k)
		}
	}
	if len(keys) > 0 {
		desc += "\nAnnotations:\n"
		sort.Sort(keys)
		for _, k := range keys {
			desc += fmt.Sprintf("%16s: %s\n",k,a.Annotations[k])
		}
	}

	keys = make(pmod.LabelNames, 0, len(a.Labels))
	for k, _ := range a.Labels {
		if ! ProcLabelKeys[k] {
			keys = append(keys, k)
		}
	}
	if len(keys) > 0 {
		desc  += "\nLabels:\n"
		sort.Sort(keys)
		for _, k := range keys {
			desc += fmt.Sprintf("%16s: %s\n",k,a.Labels[k])
		}
	}

	return desc
}

func (ag *AlertGroup) Desc() string {
	var desc string

	keys := make(pmod.LabelNames, 0, len(ag.ComAnnots))
	for k, _ := range ag.ComAnnots {
		if ! ProcAnnotKeys[k] {
			keys = append(keys, k)
		}
	}
	if len(keys) > 0 {
		desc = "\nCommon Annotations:\n"
		sort.Sort(keys)
		for _, k := range keys {
			desc += fmt.Sprintf("%16s: %s\n",k,ag.ComAnnots[k])
		}
	}

	keys = make(pmod.LabelNames, 0, len(ag.ComLabels))
	for k, _ := range ag.ComLabels {
		if ! ProcLabelKeys[k] {
			keys = append(keys, k)
		}
	}
	if len(keys) > 0 {
		desc  += "\nCommon Labels:\n"
		sort.Sort(keys)
		for _, k := range keys {
			desc += fmt.Sprintf("%16s: %s\n",k,ag.ComLabels[k])
		}
	}
	desc  += fmt.Sprintf("\nAlerts(%d):\n",len(ag.Alerts))
	for i, a := range ag.Alerts {
		desc += fmt.Sprintf("\nTitle(%d): %s\n",i+1,a.Title())
		desc += a.Desc() + "\n"
	}
	return desc
}
