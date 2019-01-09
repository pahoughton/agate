## agate

[![Test Build Status](https://travis-ci.org/pahoughton/agate.png)](https://travis-ci.org/pahoughton/agate)

Prometheus alertmanager webhook that responds to alerts by generating
tickets, executing ansible roles and running scripts.

## state

under development - see wip branch

## usage

```
usage: agate [<flags>]

prometheus alertmanager webhook processor

Flags:
  --help                   Show context-sensitive help (also try --help-long and
                           --help-man).
  --version                Show application version.
  --config-fn="agate.yml"  config filename
  --debug                  debug output to stdout
```

The webhook URL, http://hostname:port/alerts, processes alert manager
alert groups.

Prometheus metrics are available via http://hostname:port/metrics,

### config file

```yaml
listen-addr: ":1234"
ticket-default-sys: gitlab
ticket-default-grp: project
debug: true
base-dir: /var/lib/agate
max-days: 15
email-smtp: host:25
email-from: no-reply-agate@nowhere.non
gitlab-url: https://gitlab.com/api/v4
gitlab-token: secret-token
hpsm-url: https://hpsm/apiv3
hpsm-user: hpsm
hpsm-pass: pass
mock-ticket-url: http://localhost:5003/ticket
```

### labels

* ticket: gitlab|mock

  ticketing system

* gitlab: project

  gitlab project for issue creation

* hpsmwg: WGTEST

  HPSM Workgroup for ticket creation

* email: linux-queue@nowhere.none

  email address to send tickets to

* ansible: role

  execute the specified ansible role on the instance alert labels
  are added to the playbook as host vars

* script: name

  run the specified script passing instance as the

* close_resolved: bool

  close ticket when resolved

## features

Alerts generate a ticket via a ticket-url or ticket-email-to.
Comments with remediation output and resolution details are also generated.

## install

An example [sytemd service](../master/agate.service) is
provided. The service User must be able to ssh to alerting instances
with out password and have the ability to sudo for remediation.

Ticket IDs are stored by an alert key in data-dir to update tickets with
remediation results and alert resolution. Alerts older than
data-max-days are removed every 24 hours.

There is a script in the rpm directory to generate a rpm using
rpmbuild for systems that use systemd. The default data, playbook and
scripts directories are under /var/lib/agate

## validation - under developement

Execute `vagrant up` in the test directory to initialize the
validation process. This validation requires prometheus, alertmanager,
node_exporter and process-exporter be available in $GOPATH/bin.

## contribute

https://gitlab.com/pahoughton/agate

## licenses

2018-12-05 (cc) <paul4hough@gmail.com>

GNU General Public License v3.0

See [COPYING](../master/COPYING) for full text.
