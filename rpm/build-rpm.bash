#!/bin/bash
# 2019-01-03 (cc) <paul4hough@gmail.com>
#
set -x

pushd ..
go build || exit 1
popd
rm -rf BUILD
mkdir -p SPECS BUILD/usr/sbin
cp agate-el7-rpm.spec SPECS || exit 1
cp ../agate BUILD/usr/sbin || exit 1
mkdir -p BUILD/var/lib/agate/playbook/roles/service-restart/tasks
cp ../playbook/agate.yml BUILD/var/lib/agate/playbook || exit 1
cp ../playbook/roles/service-restart/tasks/main.yml \
   BUILD/var/lib/agate/playbook/roles/service-restart/tasks || exit 1
mkdir -p BUILD/var/lib/agate/scripts
cp ../scripts/disk-usage BUILD/var/lib/agate/scripts || exit 1
mkdir -p BUILD/var/lib/agate/data
mkdir -p BUILD/etc/systemd/system
cp ../agate.service BUILD/etc/systemd/system || exit 1
rpmbuild -bb --build-in-place --buildroot `pwd`/BUILD agate-el7-rpm.spec
