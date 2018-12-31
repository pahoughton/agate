/* 2018-12-21 (cc) <paul4hough@gmail.com>
   mock ticketing system
*/
package main

import (
	"net/http"
	"os"
	"path/filepath"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	log "github.com/sirupsen/logrus"

	prom  "github.com/prometheus/client_golang/prometheus"
	proma "github.com/prometheus/client_golang/prometheus/promauto"
	promh "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	app = kingpin.New(filepath.Base(os.Args[0]),
		"http dumper service")

	listenAddr = app.Flag("laddr","listen address").
		Short('l').
		Default(":5003").
		Envar("MOCK_TICKET_LISTEN").
		String()

	tickets []string

	unsupRecvd = proma.NewCounter(
		prom.CounterOpts{
			Name:      "unsupported_received_total",
			Help:      "number of unsupported requests received",
		})
	ticketRecvd = proma.NewCounter(
		prom.CounterOpts{
			Name:      "ticket_received_total",
			Help:      "number of ticket requests received",
		})
)

func main() {
	app.Version("0.0.1")
	kingpin.MustParse(app.Parse(os.Args[1:]))

	log.SetLevel(log.TraceLevel)
	log.Info("Starting ",os.Args[0],"listing to",*listenAddr)

	tickets = append(tickets,`{"title":"Dummy One","desc":"not real"}`)
	tickets = append(tickets,`{"title":"Dummy Two","desc":"not real 2"}`)

	http.Handle("/metrics", promh.Handler())
	http.HandleFunc("/list",handleList)
	http.HandleFunc("/show",handleShow)
	http.HandleFunc("/ticket",handleTicket)
	http.HandleFunc("/",handleDefault)
	log.Fatal(http.ListenAndServe(*listenAddr,nil))
}
