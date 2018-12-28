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

    yml "gopkg.in/yaml.v2"
	promp "github.com/prometheus/client_golang/prometheus"
)

func createMockApiTicket(a *AmgrAlert){

	if *args.Debug {
		fmt.Println("DEBUG: create mock api ticket for: ")
		yout, _ := yml.Marshal(*a)
		fmt.Println(string(yout))
	}

	node := strings.Split(a.Labels["instance"],":")[0]
	title := node + ": " + a.Labels["alertname"]

	desc  := "\nAnnotations:\n"
	for k, v := range a.Annotations {
		desc += k + ": " + v + "\n"
	}
	desc  += "\nLabels:\n"
	for k, v := range a.Labels {
		desc += k + ": " + v + "\n"
	}
	desc += "\nfrom: " + a.GeneratorURL + "\n"
	desc += "\nremed: " + a.RemedOut + "\n"

	tckt := map[string]string{
		"title":       title,
		"work_group":  "WGTEST",
		"start_time":  a.StartsAt.String(),
		"status":      a.Status,
		"description": desc,
	}

	tcktJson, err := json.Marshal(tckt)
	if err != nil {
		fmt.Printf("FATAL-json.Marshal: %s\n%+v\n",err.Error(),*a)
		os.Exit(2)
	}

    req, err := http.NewRequest("POST",
		*args.TicketMockURL,
		bytes.NewBuffer(tcktJson))
    if err != nil {
		fmt.Printf("FATAL-http.NewRequest-%s: %s \n%+v\n",
			*args.TicketMockURL,err.Error(),tcktJson)
		os.Exit(2)
    }

    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
		fmt.Printf("FATAL-http.Do: %s \n%+v\n",
			err.Error(),req)
		os.Exit(2)
    }

	if resp.StatusCode != 200 {
		fmt.Printf("FATAL-%s: %s \n%",
			resp.Status,err.Error())
		os.Exit(2)
	}

	prom.TicketsGend.With(
		promp.Labels{
			"type": "mock-api",
			"dest": *args.TicketMockURL}).Inc()
}
