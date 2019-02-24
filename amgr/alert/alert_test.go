/* 2019-02-15 (cc) <paul4hough@gmail.com>
   amgr/alert model validation
*/
package alert

import (
	"fmt"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	pmod "github.com/prometheus/common/model"

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

func TestNode(t *testing.T) {
	a := Alert{}
	a.Labels = pmod.LabelSet{}
	assert.Equal(t,"",a.Node())

	exp := "alert"
	expl := pmod.LabelValue(exp)
	a.Labels = pmod.LabelSet{ "agate_node": expl }
	assert.Equal(t,exp,a.Node())
	a.Labels = pmod.LabelSet{ "hostname": expl }
	assert.Equal(t,exp,a.Node())
	a.Labels = pmod.LabelSet{ "instance": pmod.LabelValue(exp+":9100") }
	assert.Equal(t,exp,a.Node())
	a.Labels = pmod.LabelSet{
		"agate_node": expl,
		"hostname": "not-exp",
	}
	assert.Equal(t,exp,a.Node())
	a.Labels = pmod.LabelSet{
		"agate_node": expl,
		"instance": "not-exp:9100",
	}
	assert.Equal(t,exp,a.Node())
	a.Labels = pmod.LabelSet{
		"agate_node": expl,
		"hostname": "notexp",
		"instance": "notexp:9100",
	}
	assert.Equal(t,exp,a.Node())
	a.Labels = pmod.LabelSet{
		"hostname": expl,
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
	exp := pmod.LabelNames{
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
	var ylabs LabelMap
	err := yaml.Unmarshal([]byte(RandLabels),&ylabs)
	assert.Nil(t,err)
	assert.Equal(t,exp,ylabs.SortedKeys())
}

func TestAlertTitle(t *testing.T) {

	atexp := "ancbed /mnt/wd4blue free 20% below 30%"
	sexp := "scbed /mnt/wd4blue free 21% below 30%"
	texp := "cbed /mnt/wd4blue free 22% below 30%"
	ta := Alert{
		pmod.Alert{
			Labels: pmod.LabelSet{
				"alertname":  "disk-usage",
				"agate_node": "ancbed",
				"hostname":   "hncbed",
				"instance":   "cbed:9100",
				"mountpoint": "/mnt/wd4blue",
			},
			Annotations: pmod.LabelSet{
				"group_title": "multiple disk free below 30%",
				"metric":      "node_filesystem_free_bytes",
				"title":       pmod.LabelValue(texp),
				"agate_title": pmod.LabelValue(atexp),
				"subject":     pmod.LabelValue(sexp),
			},
		},
		"firing",
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
	exp = string(ta.Labels["alertname"])
	assert.Equal(t,exp,ta.Title())

	delete(ta.Labels,"alertname")
	exp = "unknown"
	assert.Equal(t,exp,ta.Title())
}

func TestAlertDesc(t *testing.T) {
	startsStr := "2019-02-14T12:34:54.311358476-07:00"
	startsAt, err := time.Parse(time.RFC3339,startsStr)
	assert.Nil(t,err)
	ta := Alert{
		pmod.Alert{
			Labels: pmod.LabelSet{
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
			Annotations: pmod.LabelSet{
				"group_title": "multiple disk free below 30%",
				"metric": "node_filesystem_free_bytes",
				"agate_title": "ancbed /mnt/wd4blue free 20% below 30%",
				"subject": "scbed /mnt/wd4blue free 21% below 30%",
				"title": "cbed /mnt/wd4blue free 22% below 30%",
				"sop": "https://wiki/disk-usage",
			},
			StartsAt: startsAt,
			GeneratorURL: "http://cbed:9090/graph?g0.expr=round%28100+%2A+node_filesystem_free_bytes+%2F+node_filesystem_size_bytes%29+%3C+30\u0026g0.tab=1",
		},
		"firing",
	}
	exp := `
From: http://cbed:9090/graph?g0.expr=round%28100+%2A+node_filesystem_free_bytes+%2F+node_filesystem_size_bytes%29+%3C+30&g0.tab=1

When: 2019-02-14 12:34:54.3113 -0700

Annotations:
             sop: https://wiki/disk-usage

Labels:
       alertname: disk-usage
             app: desktop
          device: /dev/sdb1
          fstype: ext4
             job: node
        maulnode: cbed
      mountpoint: /mnt/wd4blue
            team: storage
`
	// print(ta.Desc())
	assert.Equal(t,exp,ta.Desc())
}

const GroupJson = `{
  "receiver": "agate-resolve",
  "status": "firing",
  "alerts": [
    {
      "status": "firing",
      "labels": {
        "alertname": "disk-usage",
        "app": "desktop",
        "device": "/dev/loop0",
        "fstype": "ext4",
        "instance": "cbed:9100",
        "job": "node",
        "maulnode": "cbed",
        "mongrp": "01",
        "mountpoint": "/home/paul/wip/maul/prom-poc/testdata/mnt",
        "team": "storage"
      },
      "annotations": {
        "group_title": "multiple disk free below 30%",
        "metric": "node_filesystem_free_bytes",
        "title": "cbed /home/paul/wip/maul/prom-poc/testdata/mnt free 28% below 30%"
      },
      "startsAt": "2019-02-14T12:34:54.311358476-07:00",
      "endsAt": "0001-01-01T00:00:00Z",
      "generatorURL": "http://cbed:9090/graph?g0.expr=round%28100+%2A+node_filesystem_free_bytes+%2F+node_filesystem_size_bytes%29+%3C+30\u0026g0.tab=1"
    },
    {
      "status": "firing",
      "labels": {
        "alertname": "disk-usage",
        "app": "desktop",
        "device": "/dev/nvme0n1p3",
        "fstype": "fuseblk",
        "instance": "cbed:9100",
        "job": "node",
        "maulnode": "cbed",
        "mongrp": "01",
        "mountpoint": "/media/win",
        "team": "storage"
      },
      "annotations": {
        "group_title": "multiple disk free below 30%",
        "metric": "node_filesystem_free_bytes",
        "title": "cbed /media/win free 4% below 30%"
      },
      "startsAt": "2019-02-14T12:34:54.311358476-07:00",
      "endsAt": "0001-01-01T00:00:00Z",
      "generatorURL": "http://cbed:9090/graph?g0.expr=round%28100+%2A+node_filesystem_free_bytes+%2F+node_filesystem_size_bytes%29+%3C+30\u0026g0.tab=1"
    },
    {
      "status": "firing",
      "labels": {
        "alertname": "disk-usage",
        "app": "desktop",
        "device": "/dev/sdb1",
        "fstype": "ext4",
        "instance": "cbed:9100",
        "job": "node",
        "maulnode": "cbed",
        "mongrp": "01",
        "mountpoint": "/mnt/wd4blue",
        "team": "storage"
      },
      "annotations": {
        "group_title": "multiple disk free below 30%",
        "metric": "node_filesystem_free_bytes",
        "title": "cbed /mnt/wd4blue free 22% below 30%"
      },
      "startsAt": "2019-02-14T12:34:54.311358476-07:00",
      "endsAt": "0001-01-01T00:00:00Z",
      "generatorURL": "http://cbed:9090/graph?g0.expr=round%28100+%2A+node_filesystem_free_bytes+%2F+node_filesystem_size_bytes%29+%3C+30\u0026g0.tab=1"
    }
  ],
  "groupLabels": {
    "team": "storage"
  },
  "commonLabels": {
    "alertname": "disk-usage",
    "app": "desktop",
    "instance": "cbed:9100",
    "job": "node",
    "maulnode": "cbed",
    "mongrp": "01",
    "team": "storage"
  },
  "commonAnnotations": {
    "agate_group_title": "multiple disk free below 30%",
    "metric": "node_filesystem_free_bytes"
  },
  "externalURL": "http://cbed:9093",
  "version": "4",
  "groupKey": "{}:{team=\"storage\"}"
}
`

func TestAlertGroupTitle(t *testing.T) {

	var ag AlertGroup
	assert.Nil(t,json.Unmarshal([]byte(GroupJson), &ag))
	exp := "multiple disk free below 30%"
	assert.Equal(t,exp,ag.Title())

	delete( ag.ComAnnots, "agate_group_title")
	exp = fmt.Sprintf("%d grouped alerts",len(ag.Alerts))
	assert.Equal(t,exp,ag.Title())
}

func TestAlertGroupDesc(t *testing.T) {
	var ag AlertGroup
	if err := json.Unmarshal([]byte(GroupJson), &ag); err != nil {
		t.Errorf("json.Unmarshal agrp: %s\n%v",err.Error(),GroupJson)
		return
	}
	exp := `
Common Labels:
       alertname: disk-usage
             app: desktop
             job: node
        maulnode: cbed
            team: storage

Alerts(3):

Title(1): cbed /home/paul/wip/maul/prom-poc/testdata/mnt free 28% below 30%

From: http://cbed:9090/graph?g0.expr=round%28100+%2A+node_filesystem_free_bytes+%2F+node_filesystem_size_bytes%29+%3C+30&g0.tab=1

When: 2019-02-14 12:34:54.3113 -0700

Labels:
       alertname: disk-usage
             app: desktop
          device: /dev/loop0
          fstype: ext4
             job: node
        maulnode: cbed
      mountpoint: /home/paul/wip/maul/prom-poc/testdata/mnt
            team: storage


Title(2): cbed /media/win free 4% below 30%

From: http://cbed:9090/graph?g0.expr=round%28100+%2A+node_filesystem_free_bytes+%2F+node_filesystem_size_bytes%29+%3C+30&g0.tab=1

When: 2019-02-14 12:34:54.3113 -0700

Labels:
       alertname: disk-usage
             app: desktop
          device: /dev/nvme0n1p3
          fstype: fuseblk
             job: node
        maulnode: cbed
      mountpoint: /media/win
            team: storage


Title(3): cbed /mnt/wd4blue free 22% below 30%

From: http://cbed:9090/graph?g0.expr=round%28100+%2A+node_filesystem_free_bytes+%2F+node_filesystem_size_bytes%29+%3C+30&g0.tab=1

When: 2019-02-14 12:34:54.3113 -0700

Labels:
       alertname: disk-usage
             app: desktop
          device: /dev/sdb1
          fstype: ext4
             job: node
        maulnode: cbed
      mountpoint: /mnt/wd4blue
            team: storage

`
	assert.Equal(t,exp,ag.Desc())
}

func TestAlertKey(t *testing.T) {

	title := "cbed /mnt/wd4blue free 22% below 30%"
	a := Alert{
		pmod.Alert{
			Labels: pmod.LabelSet{
				"alertname":  "disk-usage",
				"agate_node": "cbed",
				"instance":   "cbed:9100",
				"mountpoint": "/mnt/wd4blue",
			},
			Annotations: pmod.LabelSet{
				"group_title": "multiple disk free below 30%",
				"metric":      "node_filesystem_free_bytes",
				"title":       pmod.LabelValue(title),
			},
			StartsAt: time.Now(),
		},
		"firing",
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
