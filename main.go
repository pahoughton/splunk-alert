/* 2019-01-17 (cc) <paul4hough@gmail.com>
   FIXME what is this for?
*/
package main


import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/pahoughton/splunk-alert/config"
	"github.com/pahoughton/splunk-alert/alert"

	"gopkg.in/alecthomas/kingpin.v2"

	promh "github.com/prometheus/client_golang/prometheus/promhttp"
)

type CommandArgs struct {
	ConfigFn	*string
	Debug		*bool
}

func main() {

	app := kingpin.New(filepath.Base(os.Args[0]),
		"splunk alert webhook processor").
			Version("0.1.1")

	args := CommandArgs{
		ConfigFn: app.Flag("config-fn","config filename").
			Default("agate.yml").String(),
		Debug:		app.Flag("debug","debug output to stdout").
			Default("true").Bool(),
	}

	kingpin.MustParse(app.Parse(os.Args[1:]))

	fmt.Println(os.Args[0]," starting")
	fmt.Println("loading ",*args.ConfigFn)

	cfg, err := config.LoadFile(*args.ConfigFn)
	if err != nil {
		panic(err)
	}

	sphandler := alert.New(cfg,*args.Debug)

	fmt.Println("INFO: ",os.Args[0]," listening on ",cfg.Global.ListenAddr)

	amgrcnt := 0
	for _, amgr := range cfg.Amgrs {
		for _, targ := range amgr.SConfigs.Targets {
			url := fmt.Sprintf("%s://%s",amgr.Scheme, targ)
			fmt.Println("INFO: sending to ",url)
			amgrcnt += 1
		}
	}
	if amgrcnt < 1 {
		fmt.Println("FATAL: no alertmanagers configured")
		os.Exit(2)
	}
	http.Handle("/metrics",promh.Handler())
	http.Handle("/splunk/",sphandler)

	fmt.Println("FATAL: ",
		http.ListenAndServe(cfg.Global.ListenAddr,nil).
			Error())
	os.Exit(1)
}
