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
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
	pmod "github.com/prometheus/common/model"
	"github.com/pahoughton/agate/config"
	"github.com/pahoughton/agate/amgr/alert"

	"github.com/pahoughton/agate/ticket/mock"
	"github.com/pahoughton/agate/ticket/gitlab"
	"github.com/pahoughton/agate/ticket/hpsm"

)
type TSys struct {
	mock *mock.MockServer
	hpsm *hpsm.MockServer
	gitlab *gitlab.MockServer
}

func queueAGroups(t *testing.T, am *Amgr,ags ...alert.AlertGroup) {
	qlen := len(am.db.AGroupQueue())
	for i, ag := range ags {
		agjson, err := json.Marshal(ag)
		if err != nil {	t.Fatal(err) }

		req, err := http.NewRequest("POST",
			"/alerts?resolve=1",
			bytes.NewReader(agjson))
		if err != nil { t.Fatal(err) }

		rr := httptest.NewRecorder()
		assert.Equal(t,qlen + i,len(am.db.AGroupQueue()))
		am.ServeHTTP(rr,req)
		assert.Equal(t,rr.Code,http.StatusOK)
		assert.Equal(t,qlen + i + 1,len(am.db.AGroupQueue()))
	}
}

func TestRespond(t *testing.T) {
	os.Remove("testdata/data/agate.bolt")

	today := time.Now()
	cfg := config.New()
	tsys := TSys{
		hpsm: &hpsm.MockServer{},
		mock: &mock.MockServer{},
		gitlab: gitlab.NewMockServer(),
	}

	svcmock := httptest.NewServer(tsys.mock)
	defer svcmock.Close()
	cfg.Ticket.Sys.Mock.Url = svcmock.URL

	svcgitlab := httptest.NewServer(tsys.gitlab)
	defer svcgitlab.Close()
	cfg.Ticket.Sys.Gitlab.Url = svcgitlab.URL

	svchpsm := httptest.NewServer(tsys.hpsm)
	defer svchpsm.Close()
	cfg.Ticket.Sys.Hpsm.Url = svchpsm.URL

	am := New(cfg,"testdata/data",false)
	assert.NotNil(t,am)

	comlabels := pmod.LabelSet{
		"instance":		"localhost:9100",
	}
	testuniqa := 10
	testuniqag := 5
	aguniql := make([]pmod.LabelValue,0,testuniqag)
	for i := 0; i < testuniqag; i ++ {
		aguniql = append(aguniql,pmod.LabelValue(fmt.Sprintf("ag%0d",i)))
	}
	auniql := make([]pmod.LabelValue,0,testuniqa)
	for i := 0; i < testuniqa; i ++ {
		auniql = append(auniql,pmod.LabelValue(fmt.Sprintf("a-%0d",i)))
	}

	auniqt := make([]time.Time,0,testuniqa)
	for i := 0; i < testuniqa; i ++ {
		incr := time.Duration(5 * (i+1))
		auniqt = append(auniqt,today.Add(incr * time.Second))
	}
	alist := make([]alert.Alert,0,testuniqa*testuniqa)
	rlist := make([]alert.Alert,0,testuniqa*testuniqa)
	for i := 0; i < testuniqa; i ++ {
		alist = append(alist,
			alert.Alert{pmod.Alert{
				GeneratorURL:	"http://agate-nowhere/none",
				StartsAt:		auniqt[i],
				Labels:			pmod.LabelSet{
					"alertname":	"agate-remed",
					"instance":		"localhost:9100",
					"job":			"node",
					"mongrp":		"01-01",
					"component":	auniql[i],
				},},
				"firing"})
		rlist = append(rlist,
			alert.Alert{pmod.Alert{
				GeneratorURL:	"http://agate-nowhere/none",
				StartsAt:		auniqt[i],
				EndsAt:			auniqt[i].Add(5 * time.Minute),
				Labels:			pmod.LabelSet{
					"alertname":	"agate-remed",
					"instance":		"localhost:9100",
					"job":			"node",
					"mongrp":		"01-01",
					"component":	auniql[i],
				},},
				"resolved"})
	}
	aglist := make([]alert.AlertGroup,0,testuniqag)
	raglist := make([]alert.AlertGroup,0,testuniqag)
	aglist = append(aglist,
		alert.AlertGroup{
			Version:	"4",
			Receiver:	"agate-resolve",
			Status:		"firing",
			ExtURL:		"http://alertmanager",
			GroupKey:	"instance",
			ComLabels:	comlabels,
			Alerts:		alist[:1],
			GroupLabels: comlabels,
		})
	raglist = append(raglist,
		alert.AlertGroup{
			Version:	"4",
			Receiver:	"agate-resolve",
			Status:		"firing",
			ExtURL:		"http://alertmanager",
			GroupKey:	"instance",
			ComLabels:	comlabels,
			Alerts:		rlist[:1],
			GroupLabels: comlabels,
		})
	aglist = append(aglist,
		alert.AlertGroup{
			Version:	"4",
			Receiver:	"agate-resolve",
			Status:		"firing",
			ExtURL:		"http://alertmanager",
			GroupKey:	"instance",
			ComLabels:	comlabels,
			Alerts:		alist[1:4],
			GroupLabels: comlabels,
		})
	raglist = append(raglist,
		alert.AlertGroup{
			Version:	"4",
			Receiver:	"agate-resolve",
			Status:		"firing",
			ExtURL:		"http://alertmanager",
			GroupKey:	"instance",
			ComLabels:	comlabels,
			Alerts:		[]alert.Alert{ //
				alist[1],rlist[2],rlist[3],alist[4],
			},
			GroupLabels: comlabels,
		})
	raglist = append(raglist,
		alert.AlertGroup{
			Version:	"4",
			Receiver:	"agate-resolve",
			Status:		"firing",
			ExtURL:		"http://alertmanager",
			GroupKey:	"instance",
			ComLabels:	comlabels,
			Alerts:		rlist[1:4],
			GroupLabels: comlabels,
		})
 	aglist = append(aglist,
		alert.AlertGroup{
			Version:	"4",
			Receiver:	"agate-resolve",
			Status:		"firing",
			ExtURL:		"http://alertmanager",
			GroupKey:	"instance",
			ComLabels:	comlabels,
			Alerts:		alist[4:6],
			GroupLabels: comlabels,
		})
 	raglist = append(raglist,
		alert.AlertGroup{
			Version:	"4",
			Receiver:	"agate-resolve",
			Status:		"firing",
			ExtURL:		"http://alertmanager",
			GroupKey:	"instance",
			ComLabels:	comlabels,
			Alerts:		rlist[4:6],
			GroupLabels: comlabels,
		})

	assert.True(t,len(aglist) > 1)
	queueAGroups(t,am,aglist...)
	agq := am.db.AGroupQueue()
	assert.Equal(t,len(agq),len(aglist))
	mtid := tsys.mock.Tid
	for i, qid := range agq {
		assert.Equal(t,len(agq)-i,len(am.db.AGroupQueue()))
		assert.Equal(t,mtid+i,tsys.mock.Tid)
		assert.True(t,am.Respond(qid))
		assert.Equal(t,len(agq)-(i+1),len(am.db.AGroupQueue()))
		assert.Equal(t,mtid+i+1,tsys.mock.Tid)
		assert.True(t,am.Respond(qid))
		assert.Equal(t,len(agq)-(i+1),len(am.db.AGroupQueue()))
		assert.Equal(t,mtid+i+1,tsys.mock.Tid)

		for _, a := range aglist[i].Alerts {
			got := am.db.AlertGet(a.StartsAt,a.Key())
			assert.NotNil(t,got)
		}
	}
	// duplicate ags
	queueAGroups(t,am,aglist...)
	agq = am.db.AGroupQueue()
	assert.Equal(t,len(agq),len(aglist))
	mtid = tsys.mock.Tid
	thits := tsys.mock.Hits
	for i, qid := range agq {
		assert.Equal(t,len(agq)-i,len(am.db.AGroupQueue()))
		assert.True(t,am.Respond(qid))
		assert.Equal(t,len(agq)-(i+1),len(am.db.AGroupQueue()))
		assert.Equal(t,mtid,tsys.mock.Tid)
	}
	assert.Equal(t,thits,tsys.mock.Hits)

	// resolve ....
	assert.True(t,len(raglist) > 1)
	queueAGroups(t,am,raglist...)
	agq = am.db.AGroupQueue()
	assert.Equal(t,len(agq),len(raglist))
	mtid = tsys.mock.Tid
	thits = tsys.mock.Hits
	for i, qid := range agq {
		assert.Equal(t,len(agq)-i,len(am.db.AGroupQueue()))
		assert.True(t,am.Respond(qid))
		assert.Equal(t,len(agq)-(i+1),len(am.db.AGroupQueue()))
		assert.Equal(t,mtid,tsys.mock.Tid)
		assert.True(t,tsys.mock.Hits > thits)
		for _, a := range raglist[i].Alerts {
			got := am.db.AlertGet(a.StartsAt,a.Key())
			if a.Status == "resolved" {
				assert.Nil(t,got)
			} else {
				assert.NotNil(t,got)
			}
		}
	}
	for i, _ := range agq {
		for _, a := range raglist[i].Alerts {
			got := am.db.AlertGet(a.StartsAt,a.Key())
			assert.Nil(t,got)
		}
	}

	// ticket_sys: gitlab
	galist := make([]alert.Alert,0,len(alist))
	rgalist := make([]alert.Alert,0,len(alist))
	for _, a := range alist {
		a.Labels["ticket_sys"] = "gitlab"
		galist = append(galist,a)
	}
	for _, a := range rlist {
		a.Labels["ticket_sys"] = "gitlab"
		rgalist = append(rgalist,a)
	}
	tag := alert.AlertGroup{
		Version:	"4",
		Receiver:	"agate-resolve",
		Status:		"firing",
		ExtURL:		"http://alertmanager",
		GroupKey:	"instance",
		Alerts:		galist[:5],
		ComLabels:	pmod.LabelSet{
			"instance":		"localhost:9100",
			"ticket_sys":	"gitlab",
		},
		GroupLabels: pmod.LabelSet{
			"instance":		"localhost:9100",
			"ticket_sys":	"gitlab",
		},
	}
	rtag := alert.AlertGroup{
		Version:	"4",
		Receiver:	"agate-resolve",
		Status:		"firing",
		ExtURL:		"http://alertmanager",
		GroupKey:	"instance",
		Alerts:		rgalist[:5],
		ComLabels:	pmod.LabelSet{
			"instance":		"localhost:9100",
			"ticket_sys":	"gitlab",
		},
		GroupLabels: pmod.LabelSet{
			"instance":		"localhost:9100",
			"ticket_sys":	"gitlab",
		},
	}
	tid := tsys.gitlab.Tid
	thits = tsys.gitlab.Hits

	assert.Equal(t,0,len(am.db.AGroupQueue()))
	queueAGroups(t,am,tag)
	agq = am.db.AGroupQueue()
	assert.Equal(t,len(agq),1)
	assert.True(t,am.Respond(agq[0]))
	assert.Equal(t,0,len(am.db.AGroupQueue()))
	assert.Equal(t,tid + 1,tsys.gitlab.Tid)
	assert.Equal(t,thits + 1,tsys.gitlab.Hits)

	// resolve alerts
	assert.Equal(t,0,len(am.db.AGroupQueue()))
	queueAGroups(t,am,rtag)
	agq = am.db.AGroupQueue()
	assert.Equal(t,len(agq),1)
	assert.True(t,am.Respond(agq[0]))
	assert.Equal(t,0,len(am.db.AGroupQueue()))
	assert.Equal(t,tid + 1,tsys.gitlab.Tid)
	assert.True(t,thits + 1 < tsys.gitlab.Hits)

	// ticket_sys: hpsm
	halist := make([]alert.Alert,0,len(alist))
	rhalist := make([]alert.Alert,0,len(alist))
	for _, a := range alist {
		a.Labels["ticket_sys"] = "hpsm"
		halist = append(halist,a)
	}
	for _, a := range rlist {
		a.Labels["ticket_sys"] = "hpsm"
		rhalist = append(rhalist,a)
	}
	tag = alert.AlertGroup{
		Version:	"4",
		Receiver:	"agate-resolve",
		Status:		"firing",
		ExtURL:		"http://alertmanager",
		GroupKey:	"instance",
		Alerts:		[]alert.Alert{ //
			alist[1],halist[2],galist[3],halist[4],
		},
		ComLabels:	pmod.LabelSet{
			"instance":		"localhost:9100",
		},
		GroupLabels: pmod.LabelSet{
			"instance":		"localhost:9100",
		},
	}
	rtag = alert.AlertGroup{
		Version:	"4",
		Receiver:	"agate-resolve",
		Status:		"firing",
		ExtURL:		"http://alertmanager",
		GroupKey:	"instance",
		Alerts:		[]alert.Alert{ //
			rlist[1],rhalist[2],rgalist[3],rhalist[4],
		},
		ComLabels:	pmod.LabelSet{
			"instance":		"localhost:9100",
		},
		GroupLabels: pmod.LabelSet{
			"instance":		"localhost:9100",
		},
	}
	tid = tsys.hpsm.Tid
	thits = tsys.hpsm.Hits

	assert.Equal(t,0,len(am.db.AGroupQueue()))
	queueAGroups(t,am,tag)
	agq = am.db.AGroupQueue()
	assert.Equal(t,len(agq),1)
	assert.True(t,am.Respond(agq[0]))
	assert.Equal(t,0,len(am.db.AGroupQueue()))
	assert.Equal(t,tid + 1,tsys.hpsm.Tid)
	assert.Equal(t,thits + 1,tsys.hpsm.Hits)

	// resolve alerts
	assert.Equal(t,0,len(am.db.AGroupQueue()))
	queueAGroups(t,am,rtag)
	agq = am.db.AGroupQueue()
	assert.Equal(t,len(agq),1)
	assert.True(t,am.Respond(agq[0]))
	assert.Equal(t,0,len(am.db.AGroupQueue()))
	assert.Equal(t,tid + 1,tsys.gitlab.Tid)
	assert.True(t,thits + 1 < tsys.gitlab.Hits)

	exp := "agate/test"
	for i := range galist {
		galist[i].Labels["ticket_grp"] = pmod.LabelValue(exp)
		rgalist[i].Labels["ticket_grp"] = pmod.LabelValue(exp)
	}

	tag = alert.AlertGroup{
		Version:	"4",
		Receiver:	"agate-resolve",
		Status:		"firing",
		ExtURL:		"http://alertmanager",
		GroupKey:	"instance",
		Alerts:		galist[:3],
		ComLabels:	pmod.LabelSet{
			"instance":		"localhost:9100",
			"ticket_grp":	pmod.LabelValue(exp),
			"ticket_sys":	"gitlab",
		},
		GroupLabels: pmod.LabelSet{
			"instance":		"localhost:9100",
			"ticket_grp":	pmod.LabelValue(exp),
			"ticket_sys":	"gitlab",
		},
	}
	rtag = alert.AlertGroup{
		Version:	"4",
		Receiver:	"agate-resolve",
		Status:		"firing",
		ExtURL:		"http://alertmanager",
		GroupKey:	"instance",
		Alerts:		rgalist[:3],
		ComLabels:	pmod.LabelSet{
			"instance":		"localhost:9100",
			"ticket_grp":	pmod.LabelValue(exp),
			"ticket_sys":	"gitlab",
		},
		GroupLabels: pmod.LabelSet{
			"instance":		"localhost:9100",
			"ticket_grp":	pmod.LabelValue(exp),
			"ticket_sys":	"gitlab",
		},
	}

	tid = tsys.gitlab.Tid
	thits = tsys.gitlab.Hits

	assert.Equal(t,0,len(am.db.AGroupQueue()))
	queueAGroups(t,am,tag)
	agq = am.db.AGroupQueue()
	assert.Equal(t,len(agq),1)
	assert.True(t,am.Respond(agq[0]))
	assert.Equal(t,0,len(am.db.AGroupQueue()))
	assert.Equal(t,tid + 1,tsys.gitlab.Tid)
	assert.Equal(t,thits + 1,tsys.gitlab.Hits)
	assert.Equal(t,exp,tsys.gitlab.Groups[tid])

	// resolve alerts
	assert.Equal(t,0,len(am.db.AGroupQueue()))
	queueAGroups(t,am,rtag)
	agq = am.db.AGroupQueue()
	assert.Equal(t,len(agq),1)
	assert.True(t,am.Respond(agq[0]))
	assert.Equal(t,0,len(am.db.AGroupQueue()))
	assert.Equal(t,tid + 1,tsys.gitlab.Tid)
	assert.True(t,thits + 1 < tsys.gitlab.Hits)

	exp = "agate/mixtest"
	galist[0].Labels["ticket_grp"] = pmod.LabelValue(exp)
	rgalist[0].Labels["ticket_grp"] = pmod.LabelValue(exp)
	galist[2].Labels["ticket_grp"] = pmod.LabelValue(exp)
	rgalist[2].Labels["ticket_grp"] = pmod.LabelValue(exp)
	galist[3].Labels["ticket_grp"] = pmod.LabelValue(exp)
	rgalist[3].Labels["ticket_grp"] = pmod.LabelValue(exp)

	tag = alert.AlertGroup{
		Version:	"4",
		Receiver:	"agate-resolve",
		Status:		"firing",
		ExtURL:		"http://alertmanager",
		GroupKey:	"instance",
		Alerts:		galist[:5],
		ComLabels:	pmod.LabelSet{
			"instance":		"localhost:9100",
			"ticket_sys":	"gitlab",
		},
		GroupLabels: pmod.LabelSet{
			"instance":		"localhost:9100",
			"ticket_sys":	"gitlab",
		},
	}
	rtag = alert.AlertGroup{
		Version:	"4",
		Receiver:	"agate-resolve",
		Status:		"firing",
		ExtURL:		"http://alertmanager",
		GroupKey:	"instance",
		Alerts:		rgalist[:5],
		ComLabels:	pmod.LabelSet{
			"instance":		"localhost:9100",
			"ticket_sys":	"gitlab",
		},
		GroupLabels: pmod.LabelSet{
			"instance":		"localhost:9100",
			"ticket_sys":	"gitlab",
		},
	}

	tid = tsys.gitlab.Tid
	thits = tsys.gitlab.Hits

	assert.Equal(t,0,len(am.db.AGroupQueue()))
	queueAGroups(t,am,tag)
	agq = am.db.AGroupQueue()
	assert.Equal(t,len(agq),1)
	assert.True(t,am.Respond(agq[0]))
	assert.Equal(t,0,len(am.db.AGroupQueue()))
	assert.Equal(t,tid + 1,tsys.gitlab.Tid)
	assert.Equal(t,thits + 1,tsys.gitlab.Hits)
	assert.Equal(t,exp,tsys.gitlab.Groups[tid])

	// resolve alerts
	assert.Equal(t,0,len(am.db.AGroupQueue()))
	queueAGroups(t,am,rtag)
	agq = am.db.AGroupQueue()
	assert.Equal(t,len(agq),1)
	assert.True(t,am.Respond(agq[0]))
	assert.Equal(t,0,len(am.db.AGroupQueue()))
	assert.Equal(t,tid + 1,tsys.gitlab.Tid)
	assert.True(t,thits + 1 < tsys.gitlab.Hits)

	am.Close()
}
