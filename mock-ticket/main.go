/* 2018-12-21 (cc) <paul4hough@gmail.com>
   mock ticketing system
*/
package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net/http"
	"os"
	fp "path/filepath"

	"gopkg.in/alecthomas/kingpin.v2"
	"github.com/boltdb/bolt"

	promp  "github.com/prometheus/client_golang/prometheus"
	proma "github.com/prometheus/client_golang/prometheus/promauto"
	promh "github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	DBfn	= "tickets.bolt"
	Bucket	= "ticket"
)

type CommandArgs struct {
	ListenAddr	*string
	DataDir		*string
	Debug		*bool
}

type PromMetrics struct {
	Tickets		*promp.GaugeVec
	UnsupRecvd  promp.Counter
}


var (
	store *bolt.DB

	app = kingpin.New(fp.Base(os.Args[0]),
		"http dumper service").
			Version("0.0.2")

	args = CommandArgs {
		ListenAddr:	app.Flag("listen-addr","listen address").
			Default(":5001").String(),
		DataDir:	app.Flag("data-dir","ansible playbook dir").
			Default("data").String(),
		Debug:		app.Flag("debug","debug output to stdout").
			Bool(),
	}

	promNameSpace = "mock_ticket"
	prom = PromMetrics{
		Tickets: proma.NewGaugeVec(
			promp.GaugeOpts{
				Namespace: promNameSpace,
				Name:		"tickets",
				Help:		"number of tickets",
			}, []string{
				"node",
				"state",
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

	if _, err := os.Stat(*args.DataDir); err != nil {
		fmt.Println("FATAL: ",err.Error(), *args.DataDir)
		os.Exit(1)
	}

	store, err := bolt.Open(fp.Join(*args.DataDir,DBfn),0664,nil)
	if err != nil {
		fmt.Println("FATAL: open '",
			fp.Join(*args.DataDir,DBfn),"' - ",
			err.Error())
		os.Exit(1)
	}
	// set the prom gauge values from db
	err = store.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(Bucket)) // fixme skv bucket name

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {

			var t Ticket

			d := gob.NewDecoder(bytes.NewReader(v))
			if err = d.Decode(t); err != nil {
				fmt.Println("FATAL: ticket decode - ",err.Error())
				os.Exit(1)
			}

			prom.Tickets.WithLabelValues(t.Node,t.State).Inc()
		}
		return nil
	})
	if err != nil {
		fmt.Println("FATAL: db view '",err.Error())
		os.Exit(1)
	}

	http.Handle("/metrics", promh.Handler())
	http.HandleFunc("/ticket",handleTicket)
	http.HandleFunc("/list",handleList)
	http.HandleFunc("/show",handleShow)
	http.HandleFunc("/",handleDefault)

	err = http.ListenAndServe(*args.ListenAddr,nil)
	fmt.Println("FATAL: "+err.Error())
	os.Exit(1)
}
