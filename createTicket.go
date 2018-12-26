/* 2018-12-25 (cc) <paul4hough@gmail.com>
   create ticket from alert
*/
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

func createTicket(a *AmgrAlert){

	log.Debug("generating ticket for alert: ")
	fmt.Fprintf(os.Stderr,"%+v\n",*a)

	node := strings.Split(a.Labels["instance"],":")[0]
	title := node + ": " + a.Labels["alertname"]

	desc  := "\nAnnotations:\n"
	for k, v := range a.Annotations {
		desc = k + ": " + v + "\n"
	}
	desc  += "\nLabels:\n"
	for k, v := range a.Labels {
		desc = k + ": " + v + "\n"
	}
	desc += "\nfrom: " + a.GeneratorURL + "\n"

	tckt := map[string]string{
		"title":       title,
		"work_group":  "WGTEST",
		"start_time":  a.StartsAt.String(),
		"status:"      a.Status,
		"description": desc,
	}

	tcktJson, err := json.Marshal(tckt)
	if err != nil {
		log.Error(err)
		fmt.Fprintf(os.Stderr,"%+v\n",*a)
		return
	}

    req, err := http.NewRequest("POST", *ticketURL, bytes.NewBuffer(tcktJson))
    if err != nil {
        log.Error(err)
		return
    }

    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Error(err)
		return
    }

	if resp.StatusCode != 200 {
		log.Error("error sending ticket %+v\nresp: %s",a,resp.Status)
		return
	}

	ticketGend.Inc()
}
