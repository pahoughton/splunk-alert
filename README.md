## splunk-alert

[![Test Build Status](https://travis-ci.org/pahoughton/splunk-alert.png)](https://travis-ci.org/pahoughton/splunk-alert)

send alerts to alertmanager from splunk alerts and/or api searches

## Features

```yaml
---
global:
  listen_addr: ":9321"
  search_freq: 15m
  splunk_url: "http://splunk:8089/"
  splunk_user: "me"
  splunk_pass: "pass"


searchs:
  - name: splunk_openam_error
    search: "services/saved/searches/openam_errors"
    annotations:
      sop: http://wiki/sop-splunk-openam
```

## Install

Can't

## Usage

run service

## Contribute

https://github.com/pahoughton/splunk-alert

## Licenses

2019-01-16 (cc) <paul4hough@gmail.com>

[![LICENSE](http://i.creativecommons.org/l/by/4.0/80x15.png)](http://creativecommons.org/licenses/by/4.0/)
