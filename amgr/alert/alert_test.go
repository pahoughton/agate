/* 2019-02-15 (cc) <paul4hough@gmail.com>
   amgr/alert model validation
*/
package alert

import (
	"fmt"
	"encoding/json"
	"testing"
	"time"
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
func TestSortLabels(t *testing.T) {
	var ylabs LabelMap
	if err := yaml.Unmarshal([]byte(RandLabels),&ylabs); err != nil {
		t.Error(err)
	}
	skeys := ylabs.SortedKeys()
	if len(skeys) != len(SortedLabelNames) {
		t.Errorf("len: %d != %d",len(skeys),len(SortedLabelNames))
	}
	for i, v := range skeys {
		if v != SortedLabelNames[i] {
			t.Errorf("%d: %v != %v",i,v,SortedLabelNames[i])
		}
	}
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
	if ta.Title() != atexp {
		t.Errorf("title: %v != %v\n",ta.Title(),atexp)
	}
	delete(ta.Annotations,"agate_title")
	if ta.Title() != texp {
		t.Errorf("title: %v != %v\n",ta.Title(),texp)
	}
	delete(ta.Annotations,"title")
	if ta.Title() != sexp {
		t.Errorf("title: %v != %v\n",ta.Title(),sexp)
	}
	delete(ta.Annotations,"subject")
	texp = string(ta.Labels["alertname"]) + " on " + "ancbed"
	if ta.Title() != texp {
		t.Errorf("title: %v != %v\n",ta.Title(),texp)
	}
	delete(ta.Labels,"agate_node")
	texp = string(ta.Labels["alertname"]) + " on " + "hncbed"
	if ta.Title() != texp {
		t.Errorf("title: %v != %v\n",ta.Title(),texp)
	}
	delete(ta.Labels,"hostname")
	texp = string(ta.Labels["alertname"]) + " on " + "cbed"
	if ta.Title() != texp {
		t.Errorf("title: %v != %v\n",ta.Title(),texp)
	}
	delete(ta.Labels,"instance")
	texp = string(ta.Labels["alertname"])
	if ta.Title() != texp {
		t.Errorf("title: %v != %v\n",ta.Title(),texp)
	}
	delete(ta.Labels,"alertname")
	texp = "unknown"
	if ta.Title() != texp {
		t.Errorf("title: %v != %v\n",ta.Title(),texp)
	}
}

func TestAlertDesc(t *testing.T) {
	startsStr := "2019-02-14T12:34:54.311358476-07:00"
	startsAt, err := time.Parse(time.RFC3339,startsStr)
	if err != nil {
		t.Error(err)
	}
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
	texp := `
From: http://cbed:9090/graph?g0.expr=round%28100+%2A+node_filesystem_free_bytes+%2F+node_filesystem_size_bytes%29+%3C+30&g0.tab=1

When: 2019-02-14 12:34:54.311358476 -0700 MST

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
	if ta.Desc() != texp {
		t.Errorf("desc: %v\n",ta.Desc())
	}
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

	var tag AlertGroup
	if err := json.Unmarshal([]byte(GroupJson), &tag); err != nil {
		t.Errorf("json.Unmarshal agrp: %s\n%v",err.Error(),GroupJson)
		return
	}
	texp := "multiple disk free below 30%"
	if tag.Title() != texp {
		t.Errorf("group title %v != %v\n",tag.Title(),texp)
	}
	delete( tag.ComAnnots, "agate_group_title")
	texp = fmt.Sprintf("%d grouped alerts",len(tag.Alerts))
	if tag.Title() != texp {
		t.Errorf("group title %v != %v\n",tag.Title(),texp)
	}
}

func TestAlertGroupDesc(t *testing.T) {
	var tag AlertGroup
	if err := json.Unmarshal([]byte(GroupJson), &tag); err != nil {
		t.Errorf("json.Unmarshal agrp: %s\n%v",err.Error(),GroupJson)
		return
	}
	texp := `
Common Labels:
       alertname: disk-usage
             app: desktop
             job: node
        maulnode: cbed
            team: storage

Alerts(3):

Title(1): cbed /home/paul/wip/maul/prom-poc/testdata/mnt free 28% below 30%

From: http://cbed:9090/graph?g0.expr=round%28100+%2A+node_filesystem_free_bytes+%2F+node_filesystem_size_bytes%29+%3C+30&g0.tab=1

When: 2019-02-14 12:34:54.311358476 -0700 MST

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

When: 2019-02-14 12:34:54.311358476 -0700 MST

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

When: 2019-02-14 12:34:54.311358476 -0700 MST

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
	tgot := tag.Desc()
	if len(tgot) != len(texp) {
		t.Errorf("len: %d != %d\n",len(tgot),len(texp))
	}
	if tag.Desc() != texp {
		t.Errorf("desc: %v\n",tag.Desc())
	}
}
