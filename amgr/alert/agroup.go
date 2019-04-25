/* 2019-03-31 (cc) <paul4hough@gmail.com>
   FIXME what is this for?
*/

package alert

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"strconv"
	"time"
	amgrtmpl "github.com/prometheus/alertmanager/template"
)

type AlertGroup struct {
	*amgrtmpl.Data

	Version  string `json:"version"`
	GroupKey string `json:"groupKey"`
}

func (ag AlertGroup) Key() []byte {
	if len(ag.Alerts) < 1 {
		panic("no alerts in alertgroup")
	} else {
		return Alert(ag.Alerts[0]).Key()
	}
}
func (ag AlertGroup) Bytes() []byte {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	if err := enc.Encode(ag); err != nil {
		panic(err)
	}
	return b.Bytes()
}
func NewAlertGroup(b []byte) *AlertGroup {
	dec := gob.NewDecoder(bytes.NewBuffer(b))
	ag := &AlertGroup{}
	if err := dec.Decode(ag); err != nil {
		panic(err)
	}
	return ag
}
func (ag AlertGroup) StartsAt() time.Time {
	if len(ag.Alerts) < 1 {
		panic("no alerts in alertgroup")
		return time.Time{}
	} else {
		return ag.Alerts[0].StartsAt
	}
}

func (ag AlertGroup) Title() string {
	if len(ag.Alerts) == 1 {
		return Alert(ag.Alerts[0]).Title()
	} else if title, ok := ag.CommonAnnotations[GRP_TITLE_LABEL]; ok {
		return strconv.Itoa(len(ag.Alerts)) + " " + string(title)
	} else if title, ok := ag.CommonLabels[GRP_TITLE_LABEL]; ok {
		return strconv.Itoa(len(ag.Alerts)) + " " + string(title)
	} else if len(ag.CommonLabels) > 0 {
		cnt := 0
		title := ""
		for _, k := range LabelSet(ag.CommonLabels).SortedKeys() {
			cnt += 1
			if cnt > 5 {
				break;
			}
			title += ag.CommonLabels[k] + " "
		}
		return title + " " + strconv.Itoa(len(ag.Alerts)) + " alerts"
	} else {
		return strconv.Itoa(len(ag.Alerts)) + " alerts"
	}
}

func (ag *AlertGroup) Desc() string {
	var desc string

	if len(ag.Alerts) == 1 {
		return Alert(ag.Alerts[0]).Desc()
	}
	desc += "\nalertmanager: " + ag.ExternalURL + "\n"
	if len(ag.CommonLabels) > 0 {
		desc  += "common labels:\n"
		for _, k := range LabelSet(ag.CommonLabels).SortedKeys() {
			desc += fmt.Sprintf("%16s: %s\n",k,ag.CommonLabels[k])
		}
	}
	if len(ag.CommonAnnotations) > 0 {
		desc  += "\ncommon annotations:\n"
		for _, k := range LabelSet(ag.CommonAnnotations).SortedKeys() {
			desc += fmt.Sprintf("%16s: %s\n",k,ag.CommonAnnotations[k])
		}
	}

	desc  += fmt.Sprintf("\nalerts: %d\n",len(ag.Alerts))
	for i, aga := range ag.Alerts {
		a := Alert(aga)
		desc += fmt.Sprintf("\ntitle(%d): %s\n",i+1,a.Title())
		desc += a.Desc() + "\n"
	}
	return desc
}
