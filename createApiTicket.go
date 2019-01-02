/* 2018-12-25 (cc) <paul4hough@gmail.com>
   create ticket from alert
*/
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

    yml "gopkg.in/yaml.v2"
	promp "github.com/prometheus/client_golang/prometheus"
)

func createApiTicket(a *AmgrAlert) (string, error) {

	if *args.Debug {
		fmt.Println("DEBUG: create mock api ticket for: ")
		yout, _ := yml.Marshal(*a)
		fmt.Println(string(yout))
	}

	node := strings.Split(a.Labels["instance"],":")[0]
	title := node + ": " + a.Labels["alertname"]

	desc := "start_time: " + a.StartsAt.String() + "\n"

	desc += "\nAnnotations:\n"
	for k, v := range a.Annotations {
		desc += k + ": " + v + "\n"
	}
	desc  += "\nLabels:\n"
	for k, v := range a.Labels {
		desc += k + ": " + v + "\n"
	}
	desc += "\nfrom: " + a.GeneratorURL + "\n"

	tckt := map[string]string{
		"title":		title,
		"node":			node,
		"worker":		"WGTEST",
		"status":      a.Status,
		"description": desc,
	}

	tcktJson, err := json.Marshal(tckt)
	if err != nil {
		return "", fmt.Errorf("json.Marshal: %s\n%+v\n",err.Error(),*a)
	}

	resp, err := http.Post(
		*args.TicketURL,
		"application/json",
		bytes.NewReader(tcktJson))

    if err != nil {
		return "", fmt.Errorf("http.post-%s: %s \n%+v\n",
			*args.TicketURL,err.Error(),tcktJson)
    }
	defer resp.Body.Close()


	rcont, err := ioutil.ReadAll(resp.Body);
	if err != nil {
		return "", err
	}

	var rmap map[string]string

	if err := json.Unmarshal(rcont, &rmap); err != nil {
		return "", err
    }

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("resp-status: %s\n%v",resp.Status,rcont)
	}

	tid, ok := rmap["ticket"];
	if ok {
		err = adb.AddTicket(a.StartsAt,node,a.Labels["alertname"],tid)
		if err != nil {
			return "", err
		}
	} else {
		return "", fmt.Errorf("no ticket id %v",rmap)
	}

	prom.TicketsGend.With(
		promp.Labels{
			"type": "mock-api",
			"dest": *args.TicketURL}).Inc()

	return tid, nil
}
