/* 2019-02-14 (cc) <paul4hough@gmail.com>
   alertmanager alert handler
   queue for notify & remediation as needed
*/
package amgr

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	promp "github.com/prometheus/client_golang/prometheus"
	pmod "github.com/prometheus/common/model"
	"github.com/pahoughton/agate/amgr/alert"
	"github.com/pahoughton/agate/notify"
	"github.com/pahoughton/agate/notify/note"
)
const (
	TIMEFMT = "2006-01-02 15:04:05.9999 -0700"
	LBL_NSYS = "notify_sys"
	LBL_NGRP = "notify_grp"

	Url = "/api/v4/alerts"
)

func (am *Amgr)detectNSysGrp(pnsys, pgrp string, ag alert.AlertGroup ) notify.Key {

	key := notify.Key{}

	if len(pnsys) > 0 {
		key.Sys = pnsys
	} else if v, ok := ag.CommonLabels[LBL_NSYS]; ok {
		key.Sys = string(v)
	} else {
		cntmap := make(map[string]int,len(ag.Alerts))
		for _, a := range ag.Alerts {
			if v, ok := a.Labels["notify_sys"]; ok {
				cntmap[v] += 1
			}
		}
		key.Sys = am.notify.DefSys
		max := 0
		for k, c := range cntmap {
			if c > max {
				key.Sys = k
				max = c
			}
		}
	}
	if ! am.notify.ValidSys(key.Sys) {
		fmt.Printf("WARN reseting invalid nsys %s\n%v\n",pnsys,ag)
		am.metrics.errors.Inc()
		key.Sys = am.notify.DefSys
	}

	if len(pgrp) > 0 {
		key.Grp = pgrp
	} else if v, ok := ag.CommonLabels[LBL_NGRP]; ok {
		key.Grp = string(v)
	} else {
		cntmap := make(map[string]int,len(ag.Alerts))
		for _, a := range ag.Alerts {
			if v, ok := a.Labels["notify_grp"]; ok {
				cntmap[v] += 1
			}
		}
		key.Grp = am.notify.Group(key.Sys)
		max := 0
		for k, c := range cntmap {
			if c > max {
				key.Grp = k
				max = c
			}
		}
	}
	return key
}


func (am *Amgr)ServeHTTP(w http.ResponseWriter,r *http.Request) {

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		if am.debug {
			panic(err)
		} else {
			am.error(err)
		}
		return
	}
	resolve := "true"
	if len(r.FormValue("no_resolve")) > 0 {
		resolve = "false"
	}
	pnsys := r.FormValue("system")
	pgrp := r.FormValue("group")
	var ag alert.AlertGroup
	if err := json.Unmarshal(b, &ag); err != nil {
		w.WriteHeader(500)
		panic(fmt.Sprintf(
			"json.Unmarshal agrp: %s\n%v\n",
			err.Error(),b))
	}
	if len(ag.Alerts) < 1 {
		w.WriteHeader(500)
		panic("0 alerts in alertgroup")
	}
	if ag.Version != "4" {
		w.WriteHeader(500)
		panic("unsupported version")
	}

	// convert
	alerts := make([]note.Alert,len(ag.Alerts))
	remed :=  make([]note.Alert,len(ag.Alerts))
	for _, a := range ag.Alerts {

		if a.Status == "resolved" {
			continue
		}

		na := note.Alert{}
		na.Name   = a.Labels["alertname"]
		na.Starts = a.StartsAt
		na.From = a.GeneratorURL
		na.Labels = make(pmod.LabelSet,len(a.Labels))
		for k, v := range a.Labels { na.Labels[pmod.LabelName(k)] = pmod.LabelValue(v) }
		na.Labsfp = na.Labels.Fingerprint()
		alerts = append(alerts,na)

		if am.remed.HasRemed(na) {
			remed = append(remed,na)
		}
	}

	note := note.Note{
		Labels: make(pmod.LabelSet,len(ag.CommonLabels)),
		Alerts: alerts,
		From: ag.ExternalURL,
	}
	for k, v := range ag.CommonLabels { note.Labels[pmod.LabelName(k)] = pmod.LabelValue(v) }

	nkey := am.detectNSysGrp(pnsys,pgrp,ag)

	// queue notify & remed - FIXME REMEDCNT
	am.notify.Send(nkey,note,0)
	// FIXME - no remed yet
	// am.remed.Queue(nkey,alerts)

	// metrics
	ml := promp.Labels{
		"sys": nkey.Sys,
		"grp": nkey.Grp,
		"resolve":resolve,
	}
	am.metrics.groups.With(ml).Inc()
}
