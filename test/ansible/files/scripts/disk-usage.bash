#!/bin/bash
# 2019-01-08 (cc) <paul4hough@gmail.com>
#
set -x
echo args: $*
node=$1
varsfn="$2"
cat $varsfn

tmpfn=/tmp/vars.$$

# turn yaml into shell vars
grep '^[A-Za-z0-9_]*:' "$varsfn" | sed 's~\([^:]*\): *\(.*\)~\1="\2"~' > $tmpfn
cat $tmpfn
source $tmpfn

cat <<EOF | ssh $node "bash -s" --
set -x
[ -n "$mountpoint" ] || exit 1
echo $mountpoint
df -h $mountpoint

pushd $mountpoint

for dir in \`find . -type d | sort\`; do
  echo \$dir
  du -sh "\$dir"
done

EOF
