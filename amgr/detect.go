/* 2019-02-14 (cc) <paul4hough@gmail.com>
   alertmanager alert handler
   add body to data store and notify manager
*/
package amgr

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	promp "github.com/prometheus/client_golang/prometheus"
	"github.com/pahoughton/agate/db"
	"github.com/pahoughton/agate/amgr/alert"
	"github.com/pahoughton/agate/notify"
)
const (
	LBL_NSYS = "notify_sys"
	LBL_NGRP = "notify_grp"

	Url = "/api/v4/alerts"
)

func (am *Amgr)detectNSysGrp(pnsys, pgrp string, ag alert.AlertGroup ) (sys, grp string) {



	if len(pnsys) > 0 {
		sys = pnsys
	} else if v, ok := ag.CommonLabels[LBL_NSYS]; ok {
		sys = string(v)
	} else {
		cntmap := make(map[string]int,len(ag.Alerts))
		for _, a := range ag.Alerts {
			if v, ok := a.Labels["notify_sys"]; ok {
				cntmap[v] += 1
			}
		}
		sys = am.notify.DefSys()
		max := 0
		for k, c := range cntmap {
			if c > max {
				sys = k
				max = c
			}
		}
	}
	if ! notify.ValidSys(sys) {
		fmt.Printf("WARN reseting invalid nsys %s\n%v\n",pnsys,ag)
		am.metrics.errors.nsys.With(promp.Labels{
			"sys": sys,
		}).Inc()
		sys = am.notify.DefSys()
	}

	if len(pgrp) > 0 {
		grp = pgrp
	} else if v, ok := ag.CommonLabels[LBL_NGRP]; ok {
		grp = string(v)
	} else {
		cntmap := make(map[string]int,len(ag.Alerts))
		for _, a := range ag.Alerts {
			if v, ok := a.Labels["notify_grp"]; ok {
				cntmap[v] += 1
			}
		}
		grp = am.notify.Group(sys)
		max := 0
		for k, c := range cntmap {
			if c > max {
				grp = k
				max = c
			}
		}
	}
	return sys, grp
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
	nsys, ngrp = am.detectNSysGrp(pnsys,pgrp,ag)

	alerts := make([]pmod.LabelSet,len(ag.Alerts))
	remed :=  make([]pmod.LabelSet,len(ag.Alerts))
	remedCnt := 0
	for _, a := range ag.Alerts {
		als := a.Labels
		als["status"] = a.Status
		als["starts_at"] = a.StartsAt
		if len(a.EndsAt) > 0 {
			als["ends_at"] = a.EndsAt
		}
		alerts = append(alerts,als)
		if am.remed.HasRemed(als) {
			remed = append(remed,als)
		}
	}
	// queue notify & remed
	nkey := am.notify.Queue(
		sys,grp,ag.Key(),
		ag.Title(),
		ag.Desc(),
		ag.CommonLabels,
		alerts,
		remedCnt,
		resolve == "true")

	for _, a := range remed {
		am.remed.Queue(a,nkey,a["status"] == "resolved")
	}
	// metrics
	ml := promp.Labels{
		"sys": notify.NSys(nsys.Sys).String(),
		"grp": nsys.Grp,
		"resolve":resolve,
	}
	am.metrics.groups.With(ml).Inc()
}
