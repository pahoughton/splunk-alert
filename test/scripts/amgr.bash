#!/bin/bash
# 2019-01-19 (cc) <paul4hough@gmail.com>
#
set -x

scriptdir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
cfgdir=`dirname $scriptdir`/config

amgrcfg="$cfgdir/alertmanager.yml"

if [ ! -f "$amgrcfg" ] ; then
  echo "$amgrcfg not found"
  exit 1
fi

alertmanager --config.file $amgrcfg
