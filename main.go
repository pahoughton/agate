/* 2018-12-25 (cc) <paul4hough@gmail.com>
   agate entry point
*/
package main

import (
	"net/http"
	"os"
	"path/filepath"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
	log     "github.com/sirupsen/logrus"

	prom  "github.com/prometheus/client_golang/prometheus"
	proma "github.com/prometheus/client_golang/prometheus/promauto"
	promh "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	app = kingpin.New(filepath.Base(os.Args[0]),
		"prometheus alertmanager webhook processor")

	listenAddr = app.Flag("listen-addr","listen address").
		Short('l').
		Default(":5001").
		String()

	scriptDir = app.Flag("script-dir","shell script dir").
		Short('s').
		Default("scriptss").
		String()

	pbookDir = app.Flag("playbook-dir","ansible playbook dir").
		Short('p').
		Default("playbooks").
		String()

	nspace = "agate"
	alertGroupsRecvd = proma.NewCounter(
		prom.CounterOpts{
			Namespace: nspace,
			Name:      "alert_group_received_total",
			Help:      "number of alert groups received",
		})
	resolvedGroupsRecvd = proma.NewCounter(
		prom.CounterOpts{
			Namespace: nspace,
			Name:      "resolved_group_received_total",
			Help:      "number of resolved alert groups received",
		})
	alertsRecvd = proma.NewCounter(
		prom.CounterOpts{
			Namespace: nspace,
			Name:      "alert_received_total",
			Help:      "number of alerts received",
		})
	scriptProcd = proma.NewCounter(
		prom.CounterOpts{
			Namespace: nspace,
			Name:      "script_processed_total",
			Help:      "number of alerts processed with ansible",
		})
	ansibleProcd = proma.NewCounter(
		prom.CounterOpts{
			Namespace: nspace,
			Name:      "ansible_processed_total",
			Help:      "number of alerts processed with ansible",
		})
	ticketGend = proma.NewCounter(
		prom.CounterOpts{
			Namespace: nspace,
			Name:      "ticket_generated_total",
			Help:      "number of tickets generated",
		})
	unsupRecvd = proma.NewCounter(
		prom.CounterOpts{
			Namespace: nspace,
			Name:      "unsupported_received_total",
			Help:      "number of unsupported request received",
		})
)

func main() {

	app.Version("0.0.3")
	kingpin.MustParse(app.Parse(os.Args[1:]))

	log.SetLevel(log.TraceLevel)
	log.Info(os.Args[0]," started")

	if _, err := os.Stat(*pbookDir); err != nil {
		log.Fatal(err)
	}
	if _, err := os.Stat(*scriptDir); err != nil {
		log.Fatal(err)
	}
	http.Handle("/metrics", promh.Handler())
	http.HandleFunc("/alerts",handleAlertGroup)
	http.HandleFunc("/",handleUnsup)


	log.Fatal(http.ListenAndServe(*listenAddr,nil))
}
