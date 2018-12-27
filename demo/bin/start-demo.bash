#!/bin/bash
# 2018-12-26 (cc) <paul4hough@gmail.com>
#
set -x

function cleanup {
  kill $(jobs -p)
  rm -f $HOME/.ssh/id_rsa.pub
  rm -f $HOME/.ssh/id_rsa
}

trap cleanup SIGINT SIGTERM EXIT

if [ ! -f $HOME/.ssh/config ] ; then
  echo StrictHostKeyChecking no >> $HOME/.ssh/config
  chmod 600 $HOME/.ssh/config
fi

if [ ! -f $HOME/.ssh/id_rsa ] ; then
  ssh-keygen -f $HOME/.ssh/id_rsa -N ""
  cat $HOME/.ssh/id_rsa.pub >> $HOME/.ssh/authorized_keys
fi


bin/mock-service \
  --laddr ":5010" \
  --log-fn "log/mock-service.log" \
  > log/mock-service.out 2>&1 &

bin/mock-logger \
  --laddr ":5002" \
  --log-fn "log/mock-service.log" \
  > log/mock-logger.out 2>&1 &

bin/mock-ticket \
  --laddr ":5003" \
  > log/mock-ticket.out 2>&1 &

bin/node_exporter \
  --web.listen-address ":9100" \
  --collector.systemd \
  > log/node_exporter.out 2>&1 &

bin/process-exporter \
  -web.listen-address 0.0.0.0:9256 \
  -config.path config/process-exporter.yml \
  > log/process-exporter.out 2>&1 &

bin/agate \
  --listen-addr ":5001" \
  --script-dir scripts \
  --playbook-dir playbooks \
  --ticket-url http://localhost:5003/ticket \
  > log/agate.out 2>&1 &

bin/alertmanager \
  --config.file config/alertmanager.yml \
  --web.listen-address "0.0.0.0:9093" \
  --storage.path data/amgr-data \
  > log/alertmanager.out 2>&1 &

bin/prometheus \
  --config.file config/prometheus.yml \
  --web.listen-address "0.0.0.0:9090" \
  --storage.tsdb.path data/prom-data \
  > log/prometheus.out 2>&1
