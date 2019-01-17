# 2019-01-03 (cc) <paul4hough@gmail.com>
#
Name:	    agate
Version:    0.1.1
Release:    el7
Summary:    Prometheus alertmanager gateway service
License:    GPLv3
URL:	    https://github.com/pahoughton/agate

Requires(pre): /usr/sbin/useradd, /usr/bin/getent
Requires(postun): /usr/sbin/userdel

%description
Prometheus alertmanager webhook that provides remediation scripts and/or
ansible roles controled via prometheus alert labels. It generates mock
tickets and updates the comments based on remediation and alert resolution

%pre
/usr/bin/getent group %{name} || \
  /usr/sbin/groupadd -r %{name}
/usr/bin/getent passwd %{name} || \
  /usr/sbin/useradd -r -d /var/lib/%{name} -s /sbin/nologin %{name}

%postun
/usr/sbin/userdel %{name}

%files
/usr/sbin/%{name}
/etc/systemd/system/%{name}.service
%dir /var/lib/%{name}
/var/lib/%{name}/data
/var/lib/%{name}/playbook/%{name}.yml
/var/lib/%{name}/playbook/roles/service-restart/tasks/main.yml
/var/lib/%{name}/scripts/disk-usage
