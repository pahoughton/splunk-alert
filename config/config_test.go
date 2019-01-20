/* 2019-01-07 (cc) <paul4hough@gmail.com>
   config validation
*/
package config

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"gopkg.in/yaml.v2"
)

func TestLoadConfigMin(t *testing.T) {

	var cfgExp = Config{
		Global:	GlobalConfig{
			ListenAddr: ":9321",
		},
		Amgrs:	[]Amgr{
			{
				Scheme:		"https",
				SConfigs:	AmgrSConfig{
					Targets:	[]string{
						"1.2.3.4:9093",
					},
				},
			},
		},
	}

	cfgfn := "testdata/config.good.min.yml"

	cfgGot, err := LoadFile(cfgfn)

	if err != nil {
		t.Errorf("LoadFile %s: %s",cfgfn,err)
	}

	ymlGot, err := yaml.Marshal(cfgGot)
	if err != nil {
		t.Fatalf("yaml.Marshal: %s",err)
	}

	ymlExp, err := yaml.Marshal(cfgExp)
	if err != nil {
		t.Fatalf("yaml.Marshal: %s",err)
	}
	if ! reflect.DeepEqual(ymlGot, ymlExp) {
		t.Fatalf("%s: unexpected diff:\n  got:\n%s\n  exp:\n%s\n",
			cfgfn,
			ymlGot,
			ymlExp)
	}
}

func TestLoadConfigFull(t *testing.T) {

	gfreq, err := time.ParseDuration("15m")
	if err != nil {
		panic(err)
	}
	sfreq, err := time.ParseDuration("45m")
	if err != nil {
		panic(err)
	}

	var cfgExp = Config{
		Global:		GlobalConfig{
			ListenAddr: ":9321",
			Freq:		gfreq,
			SearchUrl:	"http://splunk:8089/",
			SearchUser:	"me",
			SearchPass:	"pass",
			Labels:		map[string]string{
				"source":		"splunk",
				"no_resolve":	"true",
			},
			Annots:		map[string]string{
				"extra": "stuff",
			},
		},
		Amgrs:		[]Amgr{
			{
				Scheme:		"https",
				SConfigs:	AmgrSConfig{
					Targets:	[]string{
						"1.2.3.4:9093",
						"1.2.3.5:9093",
						"1.2.3.6:9093",
					},
				},
			},
		},
		YmlAlerts:	[]YmlAlert{
			{
				Name:	"log-riak-http-access",
				Labels:	map[string]string{
					"sys":	"riak",
				},
				Annots:	map[string]string{
					"sop":		"http://wiki/sop-log-riak-http-access",
					"title":	"riak http access log error",
				},
			},
			{
				Name:	"log-riak-http-error",
				Annots:	map[string]string{
					"sop":		"http://wiki/sop-log-riak-http-error",
					"title":	"riak http error log error",
				},
			},
		},
		Searches:	[]Searches{
			{
				Name:	"log-search-stuff",
				Query:  `source status == "error"`,
				Labels:	map[string]string{
					"alab":	"bval",
				},
				Annots:	map[string]string{
					"sop":		"http://wiki/sop-log-search-stuff",
					"title":	"stuff errors",
				},
			},
			{
				Name:	"log-saved-search",
				Search: "saved-search-name",
				Freq:	sfreq,
				Url:	"http://splunkb:1234/",
				Labels:	map[string]string{
					"clab":	"dval",
				},
				Annots:	map[string]string{
					"sop": "http://wiki/sop-log-saved-search",
					"title": "stuff errors",
				},
			},
		},
		AlertMap:	map[string]AlertMap{
			"log-riak-http-access": AlertMap{
				Annots:	map[string]string{
					"sop": "http://wiki/sop-log-riak-http-access",
					"title": "riak http access log error",
				},
				Labels: map[string]string{
					"sys": "riak",
				},
			},
			"log-riak-http-error": AlertMap{
				Annots: map[string]string{
					"sop": "http://wiki/sop-log-riak-http-error",
					"title": "riak http error log error",
				},
			},
		},
	}

	cfgfn := "testdata/config.good.full.yml"

	cfgGot, err := LoadFile(cfgfn)

	if err != nil {
		t.Errorf("LoadFile %s: %s",cfgfn,err)
	}

	gotYml, err := yaml.Marshal(cfgGot)
	if err != nil {
		t.Fatalf("yaml.Marshal: %s",err)
	}

	expYml, err := yaml.Marshal(cfgExp)
	if err != nil {
		t.Fatalf("yaml.Marshal: %s",err)
	}
	gotLines := strings.Split(string(gotYml), "\n")
	expLines := strings.Split(string(expYml),"\n")

	for i, gv := range gotLines {
		if gv != expLines[i] {
			t.Fatalf("\n%s !=\n%s\nGOT:\n%s\nEXP:\n%s\n",
				gv,expLines[i],gotYml,expYml)
		}
	}
}
