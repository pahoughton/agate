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

func (am *Amgr)NewNSys(
	defsys		notify.NSys,
	pnsys, pgrp	string,
	ag			alert.AlertGroup,
	resolve		bool,
) *db.NSys {

	nsys := db.NSys{Sys: uint(defsys)}

	if len(pnsys) > 0 {
		nsys.Sys = uint(notify.NewNSys(pnsys))
	} else if v, ok := ag.CommonLabels[LBL_NSYS]; ok {
		nsys.Sys = uint(notify.NewNSys(string(v)))
	} else {
		cntmap := make(map[string]int,len(ag.Alerts))
		for _, a := range ag.Alerts {
			if v, ok := a.Labels["notify_sys"]; ok {
				cntmap[v] += 1
			}
		}
		nsys.Sys = uint(am.notify.DefSys)
		max := 0
		for k, c := range cntmap {
			if c > max {
				nsys.Sys = uint(notify.NewNSys(k))
				max = c
			}
		}
	}

	if notify.NSys(nsys.Sys) >= notify.NSysUnknown {
		fmt.Printf("WARN reseting invalid nsys %s\n%v\n",pnsys,ag)
		nsys.Sys = uint(defsys)
	}
	if len(pgrp) > 0 {
		nsys.Grp = pgrp
	} else if v, ok := ag.CommonLabels[LBL_NGRP]; ok {
		nsys.Grp = string(v)
	} else {
		cntmap := make(map[string]int,len(ag.Alerts))
		for _, a := range ag.Alerts {
			if v, ok := a.Labels["notify_grp"]; ok {
				cntmap[v] += 1
			}
		}
		nsys.Grp = am.notify.Group(notify.NSys(nsys.Sys))
		max := 0
		for k, c := range cntmap {
			if c > max {
				nsys.Grp = k
				max = c
			}
		}
	}
	nsys.Resolve = resolve
	return &nsys
}


func (am *Amgr)ServeHTTP(w http.ResponseWriter,r *http.Request) {

	fmt.Printf("url: %v\n",r.URL)
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

	agkey := ag.Key()
	nsys := am.db.AGroupNSysGet(ag.StartsAt(),agkey);
	if nsys == nil {
		nsys = am.NewNSys(am.notify.DefSys,pnsys,pgrp,ag,resolve == "true")
		am.db.AGroupQueueNSysAdd(ag.StartsAt(),*nsys,agkey,ag.Bytes())
	} else {
		am.db.AGroupQueueAdd(nsys.Sys,ag.Bytes())
	}

	ml := promp.Labels{
		"sys": notify.NSys(nsys.Sys).String(),
		"grp": nsys.Grp,
		"resolve":resolve,
	}
	am.metrics.groups.With(ml).Inc()
	am.qmgr.Notify(nsys.Sys)
}
