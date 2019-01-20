#!/bin/bash
# 2019-01-19 (cc) <paul4hough@gmail.com>
#
scriptdir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
cfgdir=`dirname $scriptdir`/config

salertcfg="$cfgdir/splunk-alert.yml"

$scriptdir/../../splunk-alert --debug --config-fn $salertcfg
