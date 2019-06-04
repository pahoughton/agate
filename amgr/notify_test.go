/* 2019-02-22 (cc) <paul4hough@gmail.com>
   FIXME what is this for?
*/
package amgr


import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	amgrtmpl "github.com/prometheus/alertmanager/template"
	"github.com/pahoughton/agate/config"
	"github.com/pahoughton/agate/db"
	"github.com/pahoughton/agate/amgr/alert"
	"github.com/pahoughton/agate/notify"
	"github.com/pahoughton/agate/notify/mock"
	"github.com/pahoughton/agate/notify/gitlab"
	"github.com/pahoughton/agate/notify/hpsm"

)
type TNSys struct {
	mock *mock.MockServer
	hpsm *hpsm.MockServer
	gitlab *gitlab.MockServer
}

func queueAGroups(t *testing.T,nsys uint,am *Amgr,dags ...amgrtmpl.Data) {
	qlen := len(am.db.AGroupQueueList(nsys))
	for i, dag := range dags {
		ag := alert.AlertGroup{
			Data: &dag,
			Version: "4",
			GroupKey: "instance",
		}
		agjson, err := json.Marshal(ag)
		if err != nil {	t.Fatal(err) }

		req, err := http.NewRequest("POST",
			"/alerts?resolve=1",
			bytes.NewReader(agjson))
		if err != nil { t.Fatal(err) }

		rr := httptest.NewRecorder()
		assert.Equal(t,qlen + i,len(am.db.AGroupQueueList(nsys)))
		am.ServeHTTP(rr,req)
		assert.Equal(t,rr.Code,http.StatusOK)
		assert.Equal(t,qlen + i + 1,len(am.db.AGroupQueueList(nsys)))
	}
}

func TestNotify(t *testing.T) {
	for _, pfn := range db.DbPrevFn {
		os.Remove(path.Join(tdir,pfn))
	}
	os.Remove(path.Join(tdir,db.DbFn))

	today, _ := time.Parse(time.RFC3339, "2019-02-20T13:12:11Z")
	cfg := config.New()
	tnsys := TNSys{
		hpsm: &hpsm.MockServer{},
		mock: &mock.MockServer{},
		gitlab: gitlab.NewMockServer(),
	}

	svcmock := httptest.NewServer(tnsys.mock)
	defer svcmock.Close()
	cfg.Notify.Sys.Mock.Url = svcmock.URL

	svcgitlab := httptest.NewServer(tnsys.gitlab)
	defer svcgitlab.Close()
	cfg.Notify.Sys.Gitlab.Url = svcgitlab.URL

	svchpsm := httptest.NewServer(tnsys.hpsm)
	defer svchpsm.Close()
	cfg.Notify.Sys.Hpsm.Url = svchpsm.URL

	am := New(cfg,"testdata/data",false)
	assert.NotNil(t,am)

	comlabels := amgrtmpl.KV{
		"instance":		"localhost:9100",
	}
	testuniqag := 5
	aguniql := make([]string,0,testuniqag)
	for i := 0; i < testuniqag; i ++ {
		aguniql = append(aguniql,fmt.Sprintf("ag%0d",i))
	}
	testuniqa := 10
	auniql := make([]string,0,testuniqa)
	for i := 0; i < testuniqa; i ++ {
		auniql = append(auniql,fmt.Sprintf("a-%0d",i))
	}

	auniqt := make([]time.Time,0,testuniqa)
	for i := 0; i < testuniqa; i ++ {
		incr := time.Duration(5 * (i+1))
		auniqt = append(auniqt,today.Add(incr * time.Second))
	}
	alist := make([]amgrtmpl.Alert,0,testuniqa)
	rlist := make([]amgrtmpl.Alert,0,testuniqa)
	for i := 0; i < testuniqa; i ++ {
		a := amgrtmpl.Alert{
			Status:			"firing",
			GeneratorURL:	"http://agate-nowhere/none",
			StartsAt:		auniqt[i],
			Labels:			amgrtmpl.KV{
				"alertname":	"agate-remed",
				"instance":		"localhost:9100",
				"job":			"node",
				"mongrp":		"01-01",
				"component":	auniql[i],
			},
		}
		alist = append(alist,a)
		a.Status = "resolved"
		a.EndsAt = a.StartsAt.Add(5 * time.Minute)
		rlist = append(rlist,a)
	}
	aglist := make([]amgrtmpl.Data,0,testuniqag)
	raglist := make([]amgrtmpl.Data,0,testuniqag)
	dag := amgrtmpl.Data{
		Receiver:		"agate-resolve",
		ExternalURL:	"http://alertmanager",
		CommonLabels:	comlabels,
		GroupLabels:	comlabels,
	}

	dag.Status = "firing"
	dag.Alerts = alist[:1]
	aglist = append(aglist,dag)

	dag.Status = "resolved"
	dag.Alerts = rlist[:1]
	raglist = append(raglist,dag)

	dag.Status = "firing"
	dag.Alerts = alist[1:4]
	aglist = append(aglist,dag)
	dag.Status = "resolved"
	dag.Alerts = []amgrtmpl.Alert{alist[1],rlist[2],rlist[3],alist[4]}
	raglist = append(raglist,dag)
	dag.Alerts = rlist[1:4]
	raglist = append(raglist,dag)

	dag.Status = "firing"
	dag.Alerts = alist[4:6]
 	aglist = append(aglist,dag)
	dag.Status = "resolved"
	dag.Alerts = rlist[4:6]
 	raglist = append(raglist,dag)

	assert.True(t,len(aglist) > 1)
	queueAGroups(t,uint(am.notify.DefSys),am,aglist...)
	agq := am.db.AGroupQueueList(uint(am.notify.DefSys))
	assert.Equal(t,len(agq),len(aglist))
	mnid := tnsys.mock.Nid
	for i, qid := range agq {
		assert.Equal(t,len(agq)-i,len(am.db.AGroupQueueList(uint(am.notify.DefSys))))
		assert.Equal(t,mnid+i,tnsys.mock.Nid)
		assert.True(t,am.Notify(am.notify.DefSys,qid))
		assert.Equal(t,len(agq)-(i+1),len(am.db.AGroupQueueList(uint(am.notify.DefSys))))
		assert.Equal(t,mnid+i+1,tnsys.mock.Nid)
		for _, a := range aglist[i].Alerts {
			got := am.db.AlertNidGet(a.StartsAt,alert.Alert(a).Key())
			assert.NotNil(t,got)
		}
	}
	// duplicate ags
	queueAGroups(t,uint(am.notify.DefSys),am,aglist...)
	agq = am.db.AGroupQueueList(uint(am.notify.DefSys))
	assert.Equal(t,len(agq),len(aglist))
	mnid = tnsys.mock.Nid
	thits := tnsys.mock.Hits
	for i, qid := range agq {
		assert.Equal(t,len(agq)-i,len(am.db.AGroupQueueList(uint(am.notify.DefSys))))
		assert.True(t,am.Notify(am.notify.DefSys,qid))
		assert.Equal(t,len(agq)-(i+1),len(am.db.AGroupQueueList(uint(am.notify.DefSys))))
		assert.Equal(t,mnid,tnsys.mock.Nid)
	}
	assert.Equal(t,thits,tnsys.mock.Hits)

	// resolve ....
	nsys := uint(am.notify.DefSys)
	assert.True(t,len(raglist) > 1)
	queueAGroups(t,nsys,am,raglist...)
	agq = am.db.AGroupQueueList(nsys)
	assert.Equal(t,len(agq),len(raglist))
	mnid = tnsys.mock.Nid
	thits = tnsys.mock.Hits
	for i, qid := range agq {
		bag := am.db.AGroupQueueGet(nsys,qid)
		require.NotNil(t,bag)
		nag := alert.NewAlertGroup(bag)
		assert.Equal(t,raglist[i],*nag.Data)
		assert.Equal(t,len(agq)-i,len(am.db.AGroupQueueList(nsys)))
		assert.True(t,am.Notify(am.notify.DefSys,qid))
		assert.Equal(t,len(agq)-(i+1),len(am.db.AGroupQueueList(nsys)))
		assert.Equal(t,mnid,tnsys.mock.Nid)
		assert.True(t,tnsys.mock.Hits > thits)

		for _, a := range raglist[i].Alerts {
			if a.Status == "resolved" {
				assert.Nil(t,
					am.db.AlertNidGet(a.StartsAt,alert.Alert(a).Key()))
			}
		}
	}
	for i, _ := range agq {
		for _, a := range raglist[i].Alerts {
			got := am.db.AlertNidGet(a.StartsAt,alert.Alert(a).Key())
			assert.Nil(t,got)
		}
	}

	// notify_sys: gitlab
	nsys = uint(notify.NSysGitlab)
	galist := make([]amgrtmpl.Alert,0,len(alist))
	rgalist := make([]amgrtmpl.Alert,0,len(alist))
	for _, a := range alist {
		a.Labels["notify_sys"] = "gitlab"
		galist = append(galist,a)
	}
	for _, a := range rlist {
		a.Labels["notify_sys"] = "gitlab"
		rgalist = append(rgalist,a)
	}
	glabs := amgrtmpl.KV{
		"instance":		"localhost:9100",
		"notify_sys":	"gitlab",
	}
	dag.CommonLabels = glabs
	dag.GroupLabels = glabs
	dag.Status = "firing"
	dag.Alerts = galist[:5]

	nid := tnsys.gitlab.Nid
	thits = tnsys.gitlab.Hits

	assert.Equal(t,0,len(am.db.AGroupQueueList(uint(nsys))))
	queueAGroups(t,uint(nsys),am,dag)
	agq = am.db.AGroupQueueList(uint(nsys))
	assert.Equal(t,len(agq),1)
	assert.True(t,am.Notify(notify.NSys(nsys),agq[0]))
	assert.Equal(t,0,len(am.db.AGroupQueueList(uint(nsys))))
	assert.Equal(t,nid + 1,tnsys.gitlab.Nid)
	assert.Equal(t,thits + 1,tnsys.gitlab.Hits)

	dag.Status = "resolved"
	dag.Alerts = rgalist[:5]

	// resolve alerts
	assert.Equal(t,0,len(am.db.AGroupQueueList(uint(nsys))))
	queueAGroups(t,uint(nsys),am,dag)
	agq = am.db.AGroupQueueList(uint(nsys))
	assert.Equal(t,len(agq),1)
	assert.True(t,am.Notify(notify.NSys(nsys),agq[0]))
	assert.Equal(t,0,len(am.db.AGroupQueueList(uint(nsys))))
	assert.Equal(t,nid + 1,tnsys.gitlab.Nid)
	assert.True(t,thits + 1 < tnsys.gitlab.Hits)

	// notify_sys: hpsm
	nsys = uint(notify.NSysHpsm)
	nid = tnsys.hpsm.Nid
	thits = tnsys.hpsm.Hits

	galist = make([]amgrtmpl.Alert,0,len(alist))
	rgalist = make([]amgrtmpl.Alert,0,len(alist))
	for _, a := range alist {
		a.Labels["notify_sys"] = "hpsm"
		galist = append(galist,a)
	}
	for _, a := range rlist {
		a.Labels["notify_sys"] = "hpsm"
		rgalist = append(rgalist,a)
	}

	dag.CommonLabels["notify_sys"] = "hpsm"
	dag.GroupLabels["notify_sys"] = "hpsm"
	dag.Alerts = galist[:5]
	dag.Status = "firing"

	assert.Equal(t,0,len(am.db.AGroupQueueList(uint(nsys))))
	queueAGroups(t,uint(nsys),am,dag)
	agq = am.db.AGroupQueueList(uint(nsys))
	assert.Equal(t,len(agq),1)
	assert.True(t,am.Notify(notify.NSys(nsys),agq[0]))
	assert.Equal(t,0,len(am.db.AGroupQueueList(uint(nsys))))
	assert.Equal(t,nid + 1,tnsys.hpsm.Nid)
	assert.Equal(t,thits + 1,tnsys.hpsm.Hits)

	// resolve alerts
	dag.Alerts = rgalist[:5]
	dag.Status = "resolved"
	assert.Equal(t,0,len(am.db.AGroupQueueList(uint(nsys))))
	queueAGroups(t,nsys,am,dag)
	agq = am.db.AGroupQueueList(nsys)
	assert.Equal(t,len(agq),1)
	assert.True(t,am.Notify(notify.NSys(nsys),agq[0]))
	assert.Equal(t,0,len(am.db.AGroupQueueList(nsys)))
	assert.Equal(t,nid + 1,tnsys.hpsm.Nid)
	assert.True(t,thits + 1 < tnsys.hpsm.Hits)


	// check notify group
	nsys = uint(notify.NSysGitlab)
	nid = tnsys.gitlab.Nid
	thits = tnsys.gitlab.Hits

	exp := "agate/test"
	for i := range galist {
		galist[i].Labels["notify_grp"] = exp
		rgalist[i].Labels["notify_grp"] = exp
	}

	glabs = amgrtmpl.KV{
		"instance":		"localhost:9100",
		"notify_grp":	exp,
		"notify_sys":	"gitlab",
	}
	dag.CommonLabels = glabs
	dag.GroupLabels = glabs
	dag.Alerts = galist[:3]
	dag.Status = "firing"

	assert.Equal(t,0,len(am.db.AGroupQueueList(uint(nsys))))
	queueAGroups(t,nsys,am,dag)
	agq = am.db.AGroupQueueList(nsys)
	assert.Equal(t,len(agq),1)
	assert.True(t,am.Notify(notify.NSys(nsys),agq[0]))
	assert.Equal(t,0,len(am.db.AGroupQueueList(nsys)))
	assert.Equal(t,nid + 1,tnsys.gitlab.Nid)
	assert.Equal(t,thits + 1,tnsys.gitlab.Hits)

	// resolve alerts
	dag.CommonLabels = glabs
	dag.GroupLabels = glabs
	dag.Alerts = rgalist[:3]
	dag.Status = "resolved"
	assert.Equal(t,0,len(am.db.AGroupQueueList(nsys)))
	queueAGroups(t,nsys,am,dag)
	agq = am.db.AGroupQueueList(nsys)
	assert.Equal(t,len(agq),1)
	assert.True(t,am.Notify(notify.NSys(nsys),agq[0]))
	assert.Equal(t,0,len(am.db.AGroupQueueList(nsys)))
	assert.Equal(t,nid + 1,tnsys.gitlab.Nid)
	assert.True(t,thits + 1 < tnsys.gitlab.Hits)

	am.Del()
}
