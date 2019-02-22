## agate

[![Test Build Status](https://travis-ci.org/pahoughton/agate.png)](https://travis-ci.org/pahoughton/agate)

Prometheus alertmanager webhook that responds to alerts by generating
tickets, executing ansible roles and running scripts.

## state

under development - see wip branch

### ci/cd phase 2

cent-vm:
COMMAND - REDIRECT!:)! - sys a [group, alert, ...] to sys b :)
COMAND Supress [alertname, blah]

alert: app-agate-alert-rate
expr: >-
agate_alert_total[1m] >  10
[5m] > 20
[10m] > 40
action: crank up innibition / summarization halt alerting

Hardening

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

### configuration

See [config.good.full.yml](../master/config/testdata/config.good.full.yml)

### metric alerts

expr: *_errors_total > 0 require ping or restart to recover

### labels

nginx push server - Why not open ftp endpoint. I'll even accept samba mount.
* ticket_sys
* ticket_grp

for use by alertmanager alert grouping

* agate_ticket_group:

 * gitlab - project for issue creation
 * hpsm - incident workgroup
 * alertmanager - receiver
 * email - email to address

### annotations

* agate_ticket_sys: gitlab|mock|hpsm|alertmanager
  * default set by config
* agate_group_title: alert group title
  * default: multiple alerts for $agate_ticket_group
* agate_title: alert level title
  * default: $alertname firing on $instance
* agate_close_resolved
* default: true

## features

Alerts generate a ticket via the 'ticket' system and ticket_group or
default to the config ticket-default-sys and
ticket-default-grp. Tickets are updated with comments that include
remediation output and resolution details.

Ticket have titles and descriptions.  The title is either the
annotation.title, annotation.subject or the labels.alertname and
labels.instance.

Duplicate alerts are logged and ignored.

A remediation ansible role and/or script will be ran when
base_dir/playbook/role/labels.alertname and/or
base_dir/scripts/labels.alertname exists.

Ticket IDs are stored by an alert key in data-dir to update tickets with
remediation results and alert resolution. Alerts older than
data-max-days are removed every 24 hours.

## install

A puppet module,
[puppet-agate](https://github.com/pahoughton/puppet-agate), and an
anisble role
[ansible-agate](https://github.com/pahoughton/ansible-agate) are
availbe for installation. The specified user must be able to ssh to
remote machines and sudo for remediation.

There is also a script in the rpm directory to generate a rpm using
rpmbuild for systems that use systemd. The default config, data, playbook and
scripts directories are created under /var/lib/agate

## validation - under developement

Execute `vagrant up` in the test directory to initialize the
validation process. This validation requires prometheus, alertmanager,
node_exporter and process-exporter be available in $GOPATH/bin.

## contribute

https://github.com/pahoughton/agate

## licenses

2018-12-05 (cc) <paul4hough@gmail.com>

GNU General Public License v3.0

See [COPYING](../master/COPYING) for full text.
