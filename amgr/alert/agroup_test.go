/* 2019-04-01 (cc) <paul4hough@gmail.com>
   FIXME what is this for?
*/
package alert

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

)

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
        "group_title": "disk free below 30%",
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
        "group_title": "disk free below 30%",
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
        "group_title": "disk free below 30%",
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
    "group_title": "disk free below 30%",
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
	exp := "3 disk free below 30%"
	assert.Equal(t,exp,ag.Title())

	delete( ag.CommonAnnotations, "group_title")
	exp = "disk-usage desktop cbed:9100 node cbed  3 alerts"
	assert.Equal(t,exp,ag.Title())
}

func TestAlertGroupDesc(t *testing.T) {
	var ag AlertGroup
	if err := json.Unmarshal([]byte(GroupJson), &ag); err != nil {
		t.Errorf("json.Unmarshal agrp: %s\n%v",err.Error(),GroupJson)
		return
	}
	exp := `
alertmanager: http://cbed:9093
common labels:
       alertname: disk-usage
             app: desktop
        instance: cbed:9100
             job: node
        maulnode: cbed
          mongrp: 01
            team: storage

common annotations:
          metric: node_filesystem_free_bytes

alerts: 3

title(1): cbed /home/paul/wip/maul/prom-poc/testdata/mnt free 28% below 30%

from: http://cbed:9090/graph?g0.expr=round%28100+%2A+node_filesystem_free_bytes+%2F+node_filesystem_size_bytes%29+%3C+30&g0.tab=1

when: 2019-02-14 12:34:54.3113 -0700

labels:
       alertname: disk-usage
             app: desktop
          device: /dev/loop0
          fstype: ext4
        instance: cbed:9100
             job: node
        maulnode: cbed
          mongrp: 01
      mountpoint: /home/paul/wip/maul/prom-poc/testdata/mnt
            team: storage

annotations:
          metric: node_filesystem_free_bytes


title(2): cbed /media/win free 4% below 30%

from: http://cbed:9090/graph?g0.expr=round%28100+%2A+node_filesystem_free_bytes+%2F+node_filesystem_size_bytes%29+%3C+30&g0.tab=1

when: 2019-02-14 12:34:54.3113 -0700

labels:
       alertname: disk-usage
             app: desktop
          device: /dev/nvme0n1p3
          fstype: fuseblk
        instance: cbed:9100
             job: node
        maulnode: cbed
          mongrp: 01
      mountpoint: /media/win
            team: storage

annotations:
          metric: node_filesystem_free_bytes


title(3): cbed /mnt/wd4blue free 22% below 30%

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

`
	// print(ag.Desc())
	assert.Equal(t,exp,ag.Desc())
}
