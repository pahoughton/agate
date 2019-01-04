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
  --help                     Show context-sensitive help (also try --help-long
                             and --help-man).
  --version                  Show application version.
  --listen-addr=":5001"      listen address
  --data-dir="data"          data dir
  --data-max-days=15         max days to keep alerts
  --playbook-dir="playbook/agate.yml"
                             ansible playbook dir
  --script-dir=SCRIPT-DIR    shell script dir
  --ticket-url=TICKET-URL    ticket service url
  --ticket-smtp=TICKET-SMTP  email ticket smtp server
  --ticket-email-to=TICKET-EMAIL-TO
                             ticket email address
  --ticket-email-from="noreply-agate@no-where.not"
                             ticket email from address
  --debug                    debug output to stdout

```

The webhook URL, http://hostname:port/alerts, processes alert manager
alert groups.

Prometheus metrics are available via http://hostname:port/metrics,

### labels

* ansible: role

  execute the specified ansible role on the instance

* ansible_vars: var=value ...

  pass variable values to the role

* script: name

  run the specified script passing instance as the

* script_arg: value

  pass the value as the scripts second argument

## features

Alerts generate a ticket via a ticket-url or ticket-email-to.
Comments with remediation output and resolution details are also generated.

## install

An example [sytemd service](../blob/master/agate.service) is
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

See [COPYING] (../blob/master/COPYING) for full text.
