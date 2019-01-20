/* 2019-01-17 (cc) <paul4hough@gmail.com>
   splunk alert handler
*/
package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"time"

	"github.com/pahoughton/splunk-alert/config"

	proma "github.com/prometheus/client_golang/prometheus/promauto"
	promp "github.com/prometheus/client_golang/prometheus"
)

type Handler struct {
	Debug		bool
	Config		*config.Config
	AlertsRecvd	*promp.CounterVec
	Errors		promp.Counter
}

type SplunkAlertResult struct {
	Source	string	`json:"sourcetype"`
	Count	uint	`json:"count"`
}

type SplunkAlert struct {
	Result	SplunkAlertResult	`json:"result"`
	Sid		string				`json:"sid"`
	Link	string				`json:"results_link"`
	Search	string				`json:"search_name"`
	Owner	string				`json:"owner"`
	App		string				`json:"app"`
}

const (
	amgrEndpoint	= "/api/v1/alerts"
	contTypeJSON	= "application/json"
)

type AmgrAlert struct {
	Labels			map[string]string	`json:"labels"`
	Annotations		map[string]string	`json:"annotations"`
	StartsAt		time.Time			`json:"startsAt"`
	GeneratorURL	string				`json:"generatorURL"`
}

func New(c *config.Config, dbg bool) *Handler {

	h := &Handler{
		Debug:	dbg,
		Config: c,
		AlertsRecvd: proma.NewCounterVec(
			promp.CounterOpts{
				Namespace: "splunk",
				Name:      "alerts_received_total",
				Help:      "number of alerts received",
			}, []string{
				"name",
			}),
		Errors: proma.NewCounter(
			promp.CounterOpts{
				Namespace: "splunk",
				Name:      "errors_total",
				Help:      "number of errors",
			}),
	}

	return h
}

func (h *Handler)ServeHTTP(w http.ResponseWriter,r *http.Request) {
	if err := h.Alert(w,r); err != nil {
		fmt.Println("ERROR: ",err.Error())
		h.Errors.Inc()
    }
}

func (h *Handler)Alert(w http.ResponseWriter,r *http.Request ) error {

	aname := path.Base(r.URL.String())

	h.AlertsRecvd.With(
		promp.Labels{
			"name": aname,
		}).Inc()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("ioutil.ReadAll - %s",err.Error())
	}
	defer r.Body.Close()

	if h.Debug {
		var dbgbuf bytes.Buffer
		if err := json.Indent(&dbgbuf, b, " ", "  "); err != nil {
			return fmt.Errorf("json.Indent: ",err.Error())
		}
		fmt.Printf("DEBUG: url: %v\n",r.URL)
		fmt.Println("DEBUG: req body\n",dbgbuf.String())
	}


	var spalert SplunkAlert
	if err := json.Unmarshal(b, &spalert); err != nil {
		return fmt.Errorf("json.Unmarshal alert: %s\n%v",err.Error(),b)
    }

	ama := AmgrAlert{
		StartsAt:		time.Now(),
		GeneratorURL:	spalert.Link,
	}
	ama.Labels = make(map[string]string)

	ama.Labels["alertname"]		= aname

	for k, v := range h.Config.Global.Labels {
		ama.Labels[k] = v
	}
	for k, v := range h.Config.AlertMap[aname].Labels {
		ama.Labels[k] = v
	}
	if len(h.Config.Global.Annots) > 0 ||
		len(h.Config.AlertMap[aname].Annots) > 0 {
		ama.Annotations = make(map[string]string)
	}

	for k, v := range h.Config.Global.Annots {
		ama.Annotations[k] = v
	}
	for k, v := range h.Config.AlertMap[aname].Annots {
		ama.Annotations[k] = v
	}

	amaList := make([]AmgrAlert,1)
	amaList[0] = ama
	amjson, err := json.Marshal(amaList)
	if err != nil {
		return err
	}

	for _, amgr := range h.Config.Amgrs {
		for _, targ := range amgr.SConfigs.Targets {
			url := fmt.Sprintf("%s://%s%s",amgr.Scheme, targ, amgrEndpoint)

			if h.Debug {
				fmt.Println("DEBUG: amgr url - ", url)
			}
			resp, err := http.Post(url,contTypeJSON,bytes.NewBuffer(amjson))
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			if resp.StatusCode != 200 {
				return fmt.Errorf("alertmanager status %s\n%s",
					resp.Status,
					string(amjson))
			} else if h.Debug {
				fmt.Printf("DEBUG: sent %s\n",string(amjson))
			}

		}
	}
	return nil
}
