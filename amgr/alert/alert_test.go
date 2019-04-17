/* 2019-02-15 (cc) <paul4hough@gmail.com>
   amgr/alert model validation
*/
package alert

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	amgrtmpl "github.com/prometheus/alertmanager/template"

	"gopkg.in/yaml.v2"
)

const (

	RandLabels = `
---
fstype: ext4
app: desktop
device: /dev/loop0
instance: cbed:9100
team: storage
job: node
alertname: disk-usage
mongrp: 01
mountpoint: /home/paul/wip/maul/prom-poc/testdata/mnt
maulnode: cbed
`
)

func TestName(t *testing.T) {
	exp := "broken"
	a := Alert{ Labels: amgrtmpl.KV{"alertname": exp } }
	assert.Equal(t,exp,a.Name())
	a = Alert{}
	assert.Equal(t,"",a.Name())
}

func TestNode(t *testing.T) {
	a := Alert{}
	a.Labels = amgrtmpl.KV{}
	assert.Equal(t,"",a.Node())

	exp := "alert"
	a.Labels = amgrtmpl.KV{ "agate_node": exp }
	assert.Equal(t,exp,a.Node())
	a.Labels = amgrtmpl.KV{ "hostname": exp }
	assert.Equal(t,exp,a.Node())
	a.Labels = amgrtmpl.KV{ "instance": exp+":9100" }
	assert.Equal(t,exp,a.Node())
	a.Labels = amgrtmpl.KV{
		"agate_node": exp,
		"hostname": "not-exp",
	}
	assert.Equal(t,exp,a.Node())
	a.Labels = amgrtmpl.KV{
		"agate_node": exp,
		"instance": "not-exp:9100",
	}
	assert.Equal(t,exp,a.Node())
	a.Labels = amgrtmpl.KV{
		"agate_node": exp,
		"hostname": "notexp",
		"instance": "notexp:9100",
	}
	assert.Equal(t,exp,a.Node())
	a.Labels = amgrtmpl.KV{
		"hostname": exp,
		"instance": "notexp:9100",
	}
	assert.Equal(t,exp,a.Node())
}


/*
var (
	SortedLabelNames = pmod.LabelNames{
		"alertname",
		"app",
		"device",
		"fstype",
		"instance",
		"job",
		"maulnode",
		"mongrp",
		"mountpoint",
		"team",
	}
)
*/
func TestSortLabels(t *testing.T) {
	exp := []string{
		"alertname",
		"app",
		"device",
		"fstype",
		"instance",
		"job",
		"maulnode",
		"mongrp",
		"mountpoint",
		"team",
	}
	var ylabs LabelSet
	err := yaml.Unmarshal([]byte(RandLabels),&ylabs)
	assert.Nil(t,err)
	assert.Equal(t,exp,ylabs.SortedKeys())
}

func TestAlertTitle(t *testing.T) {

	atexp := "ancbed /mnt/wd4blue free 20% below 30%"
	sexp := "scbed /mnt/wd4blue free 21% below 30%"
	texp := "cbed /mnt/wd4blue free 22% below 30%"
	ta := Alert{
		Labels: amgrtmpl.KV{
			"alertname":  "disk-usage",
			"agate_node": "ancbed",
			"hostname":   "hncbed",
			"instance":   "cbed:9100",
			"mountpoint": "/mnt/wd4blue",
		},
		Annotations: amgrtmpl.KV{
			"group_title": "multiple disk free below 30%",
			"metric":      "node_filesystem_free_bytes",
			"title":       texp,
			"agate_title": atexp,
			"subject":     sexp,
		},
		Status: "firing",
	}
	assert.Equal(t,atexp,ta.Title())

	delete(ta.Annotations,"agate_title")
	assert.Equal(t,texp,ta.Title())

	delete(ta.Annotations,"title")
	assert.Equal(t,sexp,ta.Title())

	delete(ta.Annotations,"subject")
	exp := string(ta.Labels["alertname"]) + " on " + "ancbed"
	assert.Equal(t,exp,ta.Title())

	delete(ta.Labels,"agate_node")
	exp = string(ta.Labels["alertname"]) + " on " + "hncbed"
	assert.Equal(t,exp,ta.Title())

	delete(ta.Labels,"hostname")
	exp = string(ta.Labels["alertname"]) + " on " + "cbed"
	assert.Equal(t,exp,ta.Title())

	delete(ta.Labels,"instance")
	exp = string(ta.Labels["alertname"]) + " on "
	assert.Equal(t,exp,ta.Title())

	delete(ta.Labels,"alertname")
	exp = "alert on "
	assert.Equal(t,exp,ta.Title())
}

func TestAlertDesc(t *testing.T) {
	startsStr := "2019-02-14T12:34:54.311358476-07:00"
	startsAt, err := time.Parse(time.RFC3339,startsStr)
	assert.Nil(t,err)
	a := Alert{
		Labels: amgrtmpl.KV{
			"alertname": "disk-usage",
			"app": "desktop",
			"device": "/dev/sdb1",
			"fstype": "ext4",
			"instance": "cbed:9100",
			"job": "node",
			"maulnode": "cbed",
			"mongrp": "01",
			"mountpoint": "/mnt/wd4blue",
			"team": "storage",
		},
		Annotations: amgrtmpl.KV{
			"group_title": "multiple disk free below 30%",
			"metric": "node_filesystem_free_bytes",
			"agate_title": "ancbed /mnt/wd4blue free 20% below 30%",
			"subject": "scbed /mnt/wd4blue free 21% below 30%",
			"title": "cbed /mnt/wd4blue free 22% below 30%",
			"sop": "https://wiki/disk-usage",
		},
		StartsAt: startsAt,
		GeneratorURL: "http://cbed:9090/graph?g0.expr=round%28100+%2A+node_filesystem_free_bytes+%2F+node_filesystem_size_bytes%29+%3C+30\u0026g0.tab=1",
		Status: "firing",
	}
	exp := `
from: http://cbed:9090/graph?g0.expr=round%28100+%2A+node_filesystem_free_bytes+%2F+node_filesystem_size_bytes%29+%3C+30&g0.tab=1

when: 2019-02-14 12:34:54.3113 -0700

labels:
       alertname: disk-usage
             app: desktop
          device: /dev/sdb1
          fstype: ext4
        instance: cbed:9100
             job: node
        maulnode: cbed
          mongrp: 01
      mountpoint: /mnt/wd4blue
            team: storage

annotations:
          metric: node_filesystem_free_bytes
             sop: https://wiki/disk-usage
`
	// print(a.Desc())
	assert.Equal(t,exp,a.Desc())
}


func TestAlertKey(t *testing.T) {

	title := "cbed /mnt/wd4blue free 22% below 30%"
	a := Alert{
		Labels: amgrtmpl.KV{
			"alertname":  "disk-usage",
			"agate_node": "cbed",
			"instance":   "cbed:9100",
			"mountpoint": "/mnt/wd4blue",
		},
		Annotations: amgrtmpl.KV{
			"group_title": "multiple disk free below 30%",
			"metric":      "node_filesystem_free_bytes",
			"title":       title,
		},
		StartsAt: time.Now(),
		Status: "firing",
	}
	exp := a.Key()
	got := a.Key()
	assert.Equal(t,exp,got)

	a.StartsAt = time.Now()
	got = a.Key()
	assert.NotEqual(t,exp,got)

	exp = got
	a.Labels["agate_node"] = "cbEd"
	got = a.Key()
	assert.NotEqual(t,exp,got)

	exp = got
	delete(a.Labels,"agate_node")
	got = a.Key()
	assert.NotEqual(t,exp,got)

	assert.Equal(t,a.Key(),a.Key())
}
