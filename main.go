/* 2018-12-25 (cc) <paul4hough@gmail.com>
   agate entry point
*/
package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"path"

	"gitlab.com/pahoughton/agate/db"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	promp "github.com/prometheus/client_golang/prometheus"
	proma "github.com/prometheus/client_golang/prometheus/promauto"
	promh "github.com/prometheus/client_golang/prometheus/promhttp"
)

type CommandArgs struct {
	ListenAddr	*string
	DataDir		*string
	DataAge		*uint
	ScriptDir	*string
	Playbook	*string
	TicketURL	*string
	SMTPAddr	*string
	EmailTo		*string
	EmailFrom	*string
	Debug		*bool
}

type PromMetrics struct {
	AlertGroupsRecvd	*promp.CounterVec
	AlertsRecvd			*promp.CounterVec
	AnsiblePlays		*promp.CounterVec
	ScriptsRun			*promp.CounterVec
	TicketsGend			*promp.CounterVec
	Errors				promp.Counter
	UnsupRecvd			promp.Counter
}

var (
	adb *db.AlertDB

	app = kingpin.New(filepath.Base(os.Args[0]),
		"prometheus alertmanager webhook processor").
			Version("0.1.1")

	args = CommandArgs{
		ListenAddr:	app.Flag("listen-addr","listen address").
			Default(":5001").String(),
		DataDir:	app.Flag("data-dir","data dir").
			Default("data").String(),
		DataAge:	app.Flag("data-max-days","max days to keep alerts").
			Default("15").Uint(),
		Playbook:	app.Flag("playbook-dir","ansible playbook dir").
			Default("playbooks").String(),
		ScriptDir:  app.Flag("script-dir","shell script dir").
			String(),
		TicketURL:	app.Flag("ticket-url","ticket service url").
			String(),
		SMTPAddr:	app.Flag("ticket-smtp","email ticket smtp server").
			String(),
		EmailTo:	app.Flag("ticket-email-to","ticket email address").
			String(),
		EmailFrom:	app.Flag("ticket-email-from","ticket email from address").
			Default("noreply-agate@no-where.not").String(),
		Debug:		app.Flag("debug","debug output to stdout").
			Default("true").Bool(),
	}

	// fixme - active alerts gauge linked to db
	promNameSpace = "agate"
	prom = PromMetrics{
		AlertGroupsRecvd: proma.NewCounterVec(
			promp.CounterOpts{
				Namespace: promNameSpace,
				Name:      "alert_group_received_total",
				Help:      "number of alert groups received",
			}, []string{
				"status",
				"receiver",
			}),
		AlertsRecvd: proma.NewCounterVec(
			promp.CounterOpts{
				Namespace: promNameSpace,
				Name:      "alerts_received_total",
				Help:      "number of alerts received",
			}, []string{
				"name",
				"node",
				"status",
			}),
		AnsiblePlays: proma.NewCounterVec(
			promp.CounterOpts{
				Namespace: promNameSpace,
				Name:      "ansible_plays_total",
				Help:      "number of ansible playbook runs",
			}, []string{
				"role",
				"status",
			}),
		ScriptsRun: proma.NewCounterVec(
			promp.CounterOpts{
				Namespace: promNameSpace,
				Name:      "script_runs_total",
				Help:      "number of script runs",
			}, []string{
				"script",
				"status",
			}),
		TicketsGend: proma.NewCounterVec(
			promp.CounterOpts{
				Namespace: promNameSpace,
				Name:      "tickets_generated_total",
				Help:      "number of ticekts created",
			}, []string{
				"type",
				"dest",
			}),
		Errors: proma.NewCounter(
			promp.CounterOpts{
				Namespace: promNameSpace,
				Name:      "errors_total",
				Help:      "number of errors",
			}),
		UnsupRecvd: proma.NewCounter(
			promp.CounterOpts{
				Namespace: promNameSpace,
				Name:      "unsupported_received_total",
				Help:      "number of unsupported request received",
			}),
	}
)

func main() {

	kingpin.MustParse(app.Parse(os.Args[1:]))

	fmt.Println(os.Args[0]," listening on ",*args.ListenAddr)

	if args.Playbook != nil {
		pbStat, err := os.Stat(*args.Playbook);
		if err != nil {
			fmt.Println("FATAL: ",*args.Playbook," - ",err.Error())
			os.Exit(1)
		}
		if pbStat.IsDir() {
			fmt.Println("FATAL: ",*args.Playbook," is dir")
			os.Exit(1)
		}
		pbDir := path.Dir(*args.Playbook)
		rDir := path.Join(pbDir,"roles")
		rStat, err := os.Stat(rDir);
		if err != nil {
			fmt.Println("FATAL: ",rDir," - ",err.Error())
			os.Exit(1)
		}
		if rStat.Mode().IsDir() != true {
			fmt.Println("FATAL: ",rDir," is not dir")
			os.Exit(1)
		}
	}

	if args.ScriptDir != nil {
		sdStat, err := os.Stat(*args.ScriptDir);
		if err != nil {
			fmt.Println("FATAL: ",*args.ScriptDir," - ",err.Error())
			os.Exit(1)
		}
		if sdStat.IsDir() != true {
			fmt.Println("FATAL: ",*args.ScriptDir," is not dir")
			os.Exit(1)
		}
	}

	var err error
	adb, err = db.Open(*args.DataDir, 0664, *args.DataAge);
	if err != nil {
		fmt.Println("FATAL: open db - ",err.Error())
		os.Exit(1)
	}


	http.Handle("/metrics",promh.Handler())
	http.Handle("/alerts",errHandler(handleAlertGroup))
	// http.HandleFunc("/",handleUnsup)

	fmt.Println("FATAL: ",http.ListenAndServe(*args.ListenAddr,nil).Error())
	os.Exit(1)
}
