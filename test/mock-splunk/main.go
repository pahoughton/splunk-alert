/* 2019-01-19 (cc) <paul4hough@gmail.com>
   mock splunk alert generator
*/
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"gopkg.in/alecthomas/kingpin.v2"
)

type CommandArgs struct {
	AlertUrl	*string
	AlertData	*string
	Debug		*bool
}

type SplunkAlertResult struct {
	Source	string	`json:"sourcetype" yaml:"sourcetype"`
	Count	uint	`json:"count" yaml:"count"`
}

type SplunkAlert struct {
	Result	SplunkAlertResult	`json:"result" yaml:"result"`
	Sid		string				`json:"sid" yaml:"sid"`
	Link	string				`json:"results_link" yaml:"results_link"`
	Search	string				`json:"search_name" yaml:"search_name"`
	Owner	string				`json:"owner" yaml:"owner"`
	App		string				`json:"app" yaml:"app"`
}

type TestAlert struct {
	Name	string				`yaml:"name"`
	Splunk	SplunkAlert			`yaml:"splunk"`
}

type TestAlertData struct {
	Alerts	[]TestAlert		`yaml:"alerts"`
}

const (
	contTypeJson	= "application/json"
)

func main() {

	app := kingpin.New(filepath.Base(os.Args[0]),
		"mock splunk alerts").
			Version("0.1.1")

	args := CommandArgs{
		AlertUrl: app.Flag("alert-url","splunk-alert url").
			Default("http://localhost:9321/splunk").String(),
		AlertData: app.Flag("alert-data","alert test data yml").
			String(),
		Debug:		app.Flag("debug","debug output to stdout").
			Default("true").Bool(),
	}

	kingpin.MustParse(app.Parse(os.Args[1:]))

	if *args.Debug {
		tad := TestAlertData{
			Alerts: []TestAlert{
				TestAlert{
					Name: "log-test",
					Splunk: SplunkAlert{
						Sid: "sid-abc",
						Link: "http://abcd",
						Search: "abc",
						Owner: "me",
						App: "search",
						Result: SplunkAlertResult{
							Source: "mongod",
							Count: 5,
						},
					},
				},
			},
		}
		tyml, err := yaml.Marshal(tad)
		if err != nil {
			panic(err)
		}

		fmt.Println(string(tyml))
	}


	fmt.Println(os.Args[0]," starting")
	fmt.Println("loading ",*args.AlertData)

	dat, err := ioutil.ReadFile(*args.AlertData)
	if err != nil {
		panic(err)
	}
	ydata := &TestAlertData{}
	err = yaml.UnmarshalStrict(dat, ydata)

	if err != nil {
		panic(err)
	}
	for _, a := range ydata.Alerts {
		s, err := json.Marshal(a.Splunk)
		if err != nil {
			panic(err)
		}
		url := *args.AlertUrl + "/" + a.Name

		if *args.Debug {
			fmt.Println("url: " + url)
			fmt.Println("name: " + a.Name)
			fmt.Println("splunk: ",string(s))
		}

		resp, err := http.Post(
			url,
			contTypeJson,
			bytes.NewBuffer(s))
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			panic(fmt.Errorf("splunk-alert status %s",resp.Status))
		}
	}
}
