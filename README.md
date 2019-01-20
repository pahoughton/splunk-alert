## splunk-alert

[![Test Build Status](https://travis-ci.org/pahoughton/splunk-alert.png)](https://travis-ci.org/pahoughton/splunk-alert)

send alerts to alertmanager from splunk alerts

planned: send alerts from splunk search results

## features

Forward alerts received via splunk webhook to prometheus alert
manager. Alert labels and annotations can be added to the alert via
global and per alert configuration. The alertname is derived from the
last component of the webhook url, i.e. the webhook url
http://host:9321/splunk/log-http-access will create an alert with the
name 'log-http-access'.

### configuration

See [config.good.full.yml](../master/config/testdata/config.good.full.yml)

## install

Can't (yet)

## usage

Run as a service on host.

Create a splunk alert with a webhook trigger. The last component of
the url will be the alertmanager alertname:
http://host:port/splunk/splunk-myservie-errors will create an alert
named splunk-myservice-errors

## contribute

https://github.com/pahoughton/splunk-alert

## licenses

2019-01-16 (cc) <paul4hough@gmail.com>

GNU General Public License v3.0

See [COPYING](../master/COPYING) for full text.
