## agate

[![Test Build Status](https://travis-ci.org/pahoughton/agate.png)](https://travis-ci.org/pahoughton/agate)

Prometheus alertmanager webhook that responds to alerts by generating
tickets, executing ansible roles and running scripts.

## Todo
- validate amgr/Manager()

## state

hardening

## usage

```
usage: agate [<flags>]

prometheus alertmanager webhook processor

Flags:
  --help                Show context-sensitive help (also try --help-long and
                        --help-man).
  --version             Show application version.
  --config="agate.yml"  config filename
  --addr=":4464"        listen address
  --data="data"         data directory
  --debug               debug output to stdout

```

The webhook URL, http://hostname:port/alerts?resolve=true, processes
alert manager alert groups. http://hostname:port/alerts defaults to false

Prometheus metrics are available via http://hostname:port/metrics,

### configuration

See [config.good.full.yml](../master/config/testdata/config.good.full.yml)

### labels

use these in your alertmanager alert grouping

* ticket_sys (mock|gitlab|hpsm) (2.0: alertmanager)
* ticket_grp: {gitlab: project, hpsm: workgroup} {alertmanager: receiver}
* group_title: default: N grouped alerts (2.0 N *title)

### annotations

* title: default: $alertname on $instance

## features

* create, update & close ticket via $ticket_sys
* execute $alertname{script|ansible} (on $instance) remediation
* hardened

## install - systemd only

puppet: [puppet-agate](https://github.com/pahoughton/puppet-agate)
ansible: [ansible-agate](https://github.com/pahoughton/puppet-agate)
docker: [docker-agate](FIXME)
source: go get github.com/pahoughton/agate

## validation - under developement

### unit: go test ./...

### system:

command: rake systest
requires: rake, ansbile, vagrant, virtualbox

## example

see [test dir](../master/test) and *_test.go

## contribute

https://github.com/pahoughton/agate

## licenses

2018-12-05 (cc) <paul4hough@gmail.com>

GNU General Public License v3.0

See [COPYING](../master/COPYING) for full text.
