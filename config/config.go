/* 2019-01-17 (cc) <paul4hough@gmail.com>
   splunk-alert configuration
*/
package config

import (
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

type AmgrSConfig struct {
	Targets	[]string	`yaml:"targets"`
}

type Amgr struct {
	Scheme		string		`yaml:"scheme"`
	SConfigs	AmgrSConfig	`yaml:"static-configs"`
}

type Searches struct {
	Name	string				`yaml:"name"`
	Query	string				`yaml:"query,omitempty"`
	Search	string				`yaml:"search,omitempty"`
	Freq	time.Duration		`yaml:"freq,omitempty"`
	Url		string				`yaml:"url,omitempty"`
	Labels	map[string]string	`yaml:"labels,omitempty"`
	Annots	map[string]string	`yaml:"annotations,omitempty"`
}

type AlertMap struct {
	Annots	map[string]string
	Labels	map[string]string
}

type YmlAlert struct {
	Name	string				`yaml:"name"`
	Labels	map[string]string	`yaml:"labels,omitempty"`
	Annots	map[string]string	`yaml:"annotations,omitempty"`
}

type GlobalConfig struct {
	ListenAddr	string				`yaml:"listen-addr"`
	Freq		time.Duration		`yaml:"search-freq,omitempty"`
	SearchUrl	string				`yaml:"splunk-url,omitempty"`
	SearchUser	string				`yaml:"splunk-user,omitempty"`
	SearchPass	string				`yaml:"splunk-pass,omitempty"`
	Labels		map[string]string	`yaml:"labels,omitempty"`
	Annots		map[string]string	`yaml:"annotations,omitempty"`
}
type Config struct {
	Global		GlobalConfig		`yaml:"global"`
	Amgrs		[]Amgr				`yaml:"alertmanagers"`
	YmlAlerts	[]YmlAlert			`yaml:"alerts,omitempty"`
	Searches	[]Searches			`yaml:"searches,omitempty"`
	AlertMap	map[string]AlertMap
}

func LoadFile(fn string) (*Config, error) {

	dat, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}
	cfg := &Config{}
	err = yaml.UnmarshalStrict(dat, cfg)
	if err != nil {
		return nil, err
	}

	cfg.AlertMap = make(map[string]AlertMap)

	for _, v := range cfg.YmlAlerts {

		cfg.AlertMap[v.Name] = AlertMap{
			Labels: make(map[string]string),
			Annots: make(map[string]string),
		}

		if len(v.Annots) > 0 {
			for ak, av := range v.Annots {
				cfg.AlertMap[v.Name].Annots[ak] = av
			}
		}
		if len(v.Labels) > 0 {
			// cfg.AlertMap[v.Name].Labels = make(map[string]string)
			for lk, lv := range v.Labels {
				cfg.AlertMap[v.Name].Labels[lk] = lv
			}
		}
	}
	return cfg, nil
}
