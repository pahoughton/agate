## agate

[![Test Build Status](https://travis-ci.org/pahoughton/agate.png)](https://travis-ci.org/pahoughton/agate)

Prometheus alertmanager webhook that responds to alerts by generating
tickets, executing ansible roles and running scripts.
## state

burn-in

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

See [config.good.full.yml](blob/master/config/testdata/config.good.full.yml)


### annotations

* group_title: alert group note (issue, ticket) title
* agate_title, title, subject: note alert title

### labels

use notify_sys and notify_grp in your alertmanager alert grouping

* notify_sys (mock|gitlab|hpsm) (2.0: alertmanager)
* notify_grp: {gitlab: project, hpsm: workgroup}

* agate_node, hostname, node, instance: remediation node
* agate_title, title, subject: note (issue, ticket) title
* group_title: note title

## api

* /metrics - prometheus metrics
* /api/v4/alerts - alertmanager v4 alertgroup

```
curl $agatehost:$port/metrics

curl -XPOST -d @alert-group.json $agatehost:$port/api/v4/alerts?system=gitlab&group=maul/alerts&no_resolve=true

```


## features

* create, update & close notes via $notify_sys
* execute $alertname{script|ansible} (on $instance) remediation
* hardened

## install - systemd only

puppet: [puppet-agate](https://github.com/pahoughton/puppet-agate)
ansible: [ansible-agate](https://github.com/pahoughton/puppet-agate)
docker: [docker-agate](FIXME)
source: go get github.com/pahoughton/agate

## Todo

- validate amgr/Manager()

## validation - under developement

### unit: go test ./...

### system (wip):

command: cd test && rake spec
requires: rake, ansbile, vagrant, virtualbox

## examples

see [test dir](blob/master/test) and *_test.go

## contribute

https://github.com/pahoughton/agate

## licenses

2018-12-05 (cc) <paul4hough@gmail.com>

GNU General Public License v3.0

See [COPYING](blob/master/COPYING) for full text.
