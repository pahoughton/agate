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
	"strconv"

	"github.com/boltdb/bolt"
)

func ticketErr(w http.ResponseWriter, desc string) error {
	errMap := map[string]string{
		"error": desc,
	}

	body, err := json.Marshal(errMap)
	w.WriteHeader(500)
	if err == nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	} else {
		body = []byte(err.Error() + " " + desc)
	}
	w.Write(body)
	return errors.New(desc)
}

func handleTicket(
	w http.ResponseWriter,
	r *http.Request ) error {

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return ticketErr(w,"ioutil.ReadAll: "+err.Error())
	}
	defer r.Body.Close()

	if *args.Debug {
		var dbgbuf bytes.Buffer
		if err := json.Indent(&dbgbuf, b, " ", "  "); err != nil {
			return ticketErr(w,"json.Indent: "+err.Error())
		}
		fmt.Printf("DEBUG req body\n%s\n",dbgbuf.String())
	}

	var rt ApiTicket
	if err := json.Unmarshal(b, &rt); err != nil {
		return ticketErr(w,"json.Unmarshal: "+err.Error()+"\n"+string(b))
    }

	if len(rt.ID) < 1 { // create ticket
		var nt Ticket
		nt.Title = rt.Title
		nt.Node = rt.Node
		nt.State = rt.State
		nt.Worker = rt.Worker
		nt.Desc = rt.Desc
		if len(rt.Comment) > 0 {
			nt.Comments = append(nt.Comments, rt.Comment)
		}
		fmt.Printf("store: %v\n",tdb.db)
		err = tdb.db.Update(func(tx *bolt.Tx) error {

			bckt := tx.Bucket([]byte(Bucket))
			if bckt == nil {
				return errors.New("bucket not found "+Bucket)
			}
			id, err := bckt.NextSequence()
			if err != nil {
				return err
			}
			rt.ID = strconv.FormatUint(id,10)

			var ntbuf bytes.Buffer
			if err = gob.NewEncoder(&ntbuf).Encode(nt); err != nil {
				return err
			}

			return bckt.Put([]byte(rt.ID),ntbuf.Bytes())
		})
		if err != nil {
			return ticketErr(w,"db update "+err.Error())
		}
		prom.Tickets.WithLabelValues(rt.Node,rt.State).Inc()
	} else {
		// update ticket
		var ut Ticket

		err = tdb.db.View(func(tx *bolt.Tx) error {

			bckt := tx.Bucket([]byte(Bucket))
			if bckt == nil {
				return errors.New("bucket not found "+Bucket)
			}

			utgob := bckt.Get([]byte(rt.ID))
			if utgob == nil {
				return errors.New("not found - "+rt.ID)
			}
			return gob.NewDecoder(bytes.NewReader(utgob)).Decode(&ut)
		})
		if err != nil {
			return ticketErr(w,"db GET "+err.Error())
		}
		prom.Tickets.WithLabelValues(ut.Node,ut.State).Dec()

		if len(rt.State) > 0 {
			ut.State = rt.State
		}
		if len(rt.Worker) > 0 {
			ut.Worker = rt.Worker
		}
		if len(rt.Comment) > 0 {
			ut.Comments = append(ut.Comments, rt.Comment)
		}

		err = tdb.db.Update(func(tx *bolt.Tx) error {

			bckt := tx.Bucket([]byte(Bucket))
			if bckt == nil {
				return errors.New("bucket not found "+Bucket)
			}

			if err := bckt.Delete([]byte(rt.ID)); err != nil {
				return err
			}

			var utbuf bytes.Buffer
			if err = gob.NewEncoder(&utbuf).Encode(ut); err != nil {
				return err
			}
			return bckt.Put([]byte(rt.ID),utbuf.Bytes())
		})
		if err != nil {
			return ticketErr(w,"update "+err.Error())
		}
		prom.Tickets.WithLabelValues(ut.Node,ut.State).Inc()
	}
	tidMap := map[string]string{
		"id": rt.ID,
	}
	tidBody, err := json.Marshal(tidMap)
	w.WriteHeader(200)
	if err == nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	} else {
		tidBody = []byte(err.Error())
	}
	w.Write(tidBody)
	return err
}
