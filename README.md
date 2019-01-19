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

### annotations

* ticket: gitlab|mock

 * ticketing system

* ticket_group:

 * gitlab - project for issue creation
 * hpsm - incident workgroup
 * email - email to address

* subject: name

 * ticket title (secondary)

* title: name

 * ticket title (primary)

* close_resolved: bool

 * close ticket when resolved

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
