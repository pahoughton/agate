#!/bin/bash
# 2019-02-11 (cc) <paul4hough@gmail.com>
#
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
cd $dir
for a in alertmanager node_exporter prometheus promtool; do
  rm -f $a
  ln -s `which $a`
done
