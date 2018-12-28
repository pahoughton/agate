/* 2018-12-25 (cc) <paul4hough@gmail.com>
   agate entry point
*/
package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	promp "github.com/prometheus/client_golang/prometheus"
	proma "github.com/prometheus/client_golang/prometheus/promauto"
	promh "github.com/prometheus/client_golang/prometheus/promhttp"
)

type CommandArgs struct {
	ListenAddr		*string
	ScriptDir		*string
	PlaybookDir		*string
	TicketMockURL   *string
	SMTPAddr		*string
	EmailTo			*string
	EmailFrom		*string
	Debug			*bool
}

type PromMetrics struct {
	AlertGroupsRecvd *promp.CounterVec
	AlertsRecvd      *promp.CounterVec
	AnsiblePlays     *promp.CounterVec
	ScriptsRun       *promp.CounterVec
	TicketsGend      *promp.CounterVec
	UnsupRecvd       promp.Counter
}

var (
	app = kingpin.New(filepath.Base(os.Args[0]),
		"prometheus alertmanager webhook processor").
			Version("0.1.1")

	args = CommandArgs{
		ListenAddr:	app.Flag("listen-addr","listen address").
			Default(":5001").String(),
		PlaybookDir: app.Flag("playbook-dir","ansible playbook dir").
			Default("playbooks").String(),
		ScriptDir:  app.Flag("script-dir","shell script dir").
			String(),
		TicketMockURL:	app.Flag("ticket-mock-url","mock ticket service url").
			String(),
		SMTPAddr:	app.Flag("ticket-smtp","email ticket smtp server").
			String(),
		EmailTo:	app.Flag("ticket-email-to","ticket email address").
			String(),
		EmailFrom:	app.Flag("ticket-email-from","ticket email from address").
			Default("noreply-agate@no-where.not").String(),
		Debug:		app.Flag("debug","debug output to stdout").
			Bool(),
	}

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
				"playbook",
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
		TicketsGend:  proma.NewCounterVec(
			promp.CounterOpts{
				Namespace: promNameSpace,
				Name:      "tickets_generated_total",
				Help:      "number of ticekts created",
			}, []string{
				"type",
				"dest",
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

	if _, err := os.Stat(*args.PlaybookDir); err != nil {
		fmt.Println("FATAL: ",err.Error(), *args.PlaybookDir)
	}
	if args.ScriptDir != nil {
		if _, err := os.Stat(*args.ScriptDir); err != nil {
			fmt.Println("FATAL: ",err.Error(), *args.ScriptDir)
		}
	}
	http.Handle("/metrics", promh.Handler())
	http.HandleFunc("/alerts",handleAlertGroup)
	http.HandleFunc("/",handleUnsup)

	fmt.Println("ERROR: ",
		http.ListenAndServe(*args.ListenAddr,nil).Error())
	os.Exit(1)
}
