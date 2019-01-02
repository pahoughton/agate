/* 2018-12-21 (cc) <paul4hough@gmail.com>
   handle ticket http requests
*/
package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/boltdb/bolt"
)

func ticketRespond(w http.ResponseWriter, status int, body string) {
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write([]byte(body))
}

func handleTicket(
	w http.ResponseWriter,
	r *http.Request ) {

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("FATAL-ioutil.ReadAll: %s",err.Error())
		ticketRespond(w,500,`{"error":"read body - `+err.Error()+`"}`)
		os.Exit(2)
	}
	defer r.Body.Close()

	if *args.Debug {
		var dbgbuf bytes.Buffer
		if err := json.Indent(&dbgbuf, b, " >", "  "); err != nil {
			fmt.Println("FATAL-json.Indent: ",err.Error())
			ticketRespond(w,500,`{"error":"json indent - `+err.Error()+`"}`)
			os.Exit(2)
		}
		fmt.Printf("DEBUG req body\n%s\n",dbgbuf.String())
	}

	var rt ApiTicket
	if err := json.Unmarshal(b, &rt); err != nil {
		fmt.Println("FATAL-json.Unmarshal: %s\n%v",err.Error(),b)
		ticketRespond(w,500,`{"error":"json unmarshal - `+err.Error()+`"}`)
		return
    }

	if len(rt.ID) < 1 { // create ticket
		var nt Ticket
		nt.Title = rt.Title
		nt.Node = rt.Node
		nt.State = rt.State
		nt.Worker = rt.Worker
		nt.Desc = rt.Desc
		nt.Comments = append(nt.Comments, rt.Comment)

		err = store.Update(func(tx *bolt.Tx) error {

			bckt := tx.Bucket([]byte(Bucket))

			id, _ := bckt.NextSequence()
			rt.ID = strconv.FormatUint(id,10)

			var ntbuf bytes.Buffer
			gob.NewEncoder(&ntbuf).Encode(nt)

			return bckt.Put([]byte(rt.ID),ntbuf.Bytes())
		})
		if err != nil {
			fmt.Println("ERROR: db update ",err.Error())
			ticketRespond(w,500,`{"error":"db update - `+err.Error()+`"}`)
			return
		}
		prom.Tickets.WithLabelValues(rt.Node,rt.State).Inc()
	} else {
		// update ticket
		var ut Ticket

		err = store.View(func(tx *bolt.Tx) error {
			bckt := tx.Bucket([]byte(Bucket))
			utgob := bckt.Get([]byte(rt.ID))
			if utgob == nil {
				return errors.New("not found - "+rt.ID)
			}
			return gob.NewDecoder(bytes.NewReader(utgob)).Decode(ut)
		})
		if err != nil {
			fmt.Println("ERROR: db GET ",err.Error())
			ticketRespond(w,500,`{"error":"db get - `+err.Error()+`"}`)
			return
		}
		prom.Tickets.WithLabelValues(ut.Node,ut.State).Dec()

		if len(rt.State) > 0 {
			ut.State = rt.State
		}
		if len(rt.Worker) > 0 {
			ut.Worker = rt.Worker
		}
		ut.Comments = append(ut.Comments, rt.Comment)

		err = store.Update(func(tx *bolt.Tx) error {

			bckt := tx.Bucket([]byte(Bucket))
			if err := bckt.Delete([]byte(rt.ID)); err != nil {
				return err
			}

			var utbuf bytes.Buffer
			gob.NewEncoder(&utbuf).Encode(ut)
			return bckt.Put([]byte(rt.ID),utbuf.Bytes())
		})
		if err != nil {
			fmt.Println("ERROR: update ",err.Error())
			ticketRespond(w,500,`{"error":"update - `+err.Error()+`"}`)
			return
		}
		prom.Tickets.WithLabelValues(ut.Node,ut.State).Inc()
	}
	ticketRespond(w,200,`{"ticket":"`+rt.ID+`"}"`)
}
