 /* 2018-12-25 (cc) <paul4hough@gmail.com>
   agate entry point
*/
package main

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"runtime"

	"github.com/pahoughton/agate/config"
	"github.com/pahoughton/agate/amgr"

	"gopkg.in/alecthomas/kingpin.v2"

	promh "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	Version	  string
	Revision  string
	Branch    string
	BuildUser string
	BuildDate string
	GoVersion = runtime.Version()
)

type CommandArgs struct {
	ConfigFn	*string
	Listen		*string
	DataDir		*string
	Debug		*bool
}

func main() {
	version := fmt.Sprintf(`%s: version %s branch: %s, rev: %s
  build: %s %s
`,
		path.Base(os.Args[0]),
		Version,
		Branch,
		Revision,
		BuildDate,
		GoVersion)

	app := kingpin.New(path.Base(os.Args[0]),
		"prometheus alertmanager webhook processor").
			Version(version)

	args := CommandArgs{
		ConfigFn:	app.Flag("config","config filename").
			Default("agate.yml").String(),
		Listen:		app.Flag("addr","listen address").
			Default(":4464").String(),
		DataDir:	app.Flag("data","data directory").
			Default("data").String(),
		Debug:		app.Flag("debug","debug output to stdout").Bool(),
	}

	debug := false
	if args.Debug != nil && *args.Debug {
		debug = true
	}
	kingpin.MustParse(app.Parse(os.Args[1:]))

	fmt.Println(os.Args[0]," starting")
	fmt.Println("loading ",*args.ConfigFn)

	cfg, err := config.Load(*args.ConfigFn)
	if err != nil {
		panic(err)
	}

	if *args.Debug {
		os.Setenv("DEBUG","true")
	}

	am := amgr.New(cfg,*args.DataDir,debug)

	fmt.Println("listening on ",*args.Listen)
	am.Manage()

	http.Handle("/metrics",promh.Handler())
	http.Handle(amgr.Url,am)

	fmt.Println("FATAL: ",http.ListenAndServe(*args.Listen,nil).Error())
	os.Exit(1)
}
