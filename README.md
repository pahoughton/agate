## agate

[![Test Build Status](https://travis-ci.org/pahoughton/agate.png)](https://travis-ci.org/pahoughton/agate)

prometheus alertmanager alert gateway

## state

under development - see wip branch for latest push

## features

run ansible playbook for alerts with an ansible label. The label's
value will be the playbook run.

```json
{
    "alerts": [
        {
            "annotations": {
                "description": "thunderbird is down"
            },
            "endsAt": "2018-12-23T03:16:46.280924258-07:00",
            "generatorURL": "http://cbed:9090/graph?g0.expr=namedprocess_namegroup_num_procs%7Bgroupname%3D%22thunderbird%22%2Cjob%3D%22proc%22%7D+%3D%3D+0&g0.tab=1",
            "labels": {
                "alertname": "thunderbird-down",
                "ansible": "thunderbird-restart",
                "groupname": "thunderbird",
                "instance": "localhost:9256",
                "job": "proc"
            },
            "startsAt": "2018-12-23T03:16:31.280924258-07:00",
            "status": "firing"
        }
    ],
    "receiver": "ansible",
    "status": "firing",
    "version": "4"
}
```

will run

```
ansible-playbook -i /tmp/tmpinvXXX playbooks/thunderbird-restart.yml
```

## install

Can't

## usage

You wouldn't want to.

## contribute

https://gitlab.com/pahoughton/agate

## licenses

2018-12-25 (cc) <paul4hough@gmail.com>

GNU General Public License v3.0

See `COPYING <COPYING>`_ to see the full text.
