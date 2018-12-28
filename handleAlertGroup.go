/* 2018-12-25 (cc) <paul4hough@gmail.com>
   alert group http handler
*/
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	promp "github.com/prometheus/client_golang/prometheus"
)

func handleAlertGroup(
	w http.ResponseWriter,
	r *http.Request ) {

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("FATAL-ioutil.ReadAll: %s",err.Error())
		os.Exit(2)
	}
	defer r.Body.Close()

	if *args.Debug {
		var dbgbuf bytes.Buffer
		if err := json.Indent(&dbgbuf, b, " >", "  "); err != nil {
			fmt.Println("FATAL-json.Indent: ",err.Error())
			os.Exit(2)
		}
		fmt.Printf("DEBUG req body\n%s\n",dbgbuf.String())
	}

	var abody AmgrAlertBody
	if err := json.Unmarshal(b, &abody); err != nil {
		fmt.Println("FATAL-json.Unmarshal: %s\n%v",err.Error(),b)
		os.Exit(2)
    }

	if abody.Version != "4" {
		fmt.Println("FATAL-version: %s",abody.Version)
		os.Exit(2)
    }
	prom.AlertGroupsRecvd.With(
		promp.Labels{
			"status": abody.Status,
			"receiver": abody.Receiver,
		}).Inc()

	// ignore resolved
	if abody.Status == "resolved" {
		return
	}
	if abody.Status != "firing" {
		fmt.Println("FATAL-status: %s unsupported",abody.Status)
		os.Exit(2)
	}
	for _, alert := range abody.Alerts {
		node := strings.Split(alert.Labels["instance"],":")[0]
		prom.AlertsRecvd.With(
			promp.Labels{
				"name": alert.Labels["alertname"],
				"node": node,
				"status": abody.Status,
			}).Inc()
		if _, ok := alert.Labels["ansible"]; ok {
			procAnsible(&alert)
		} else if _, ok := alert.Labels["script"]; ok {
			procScript(&alert)
		} else {
			createTicket(&alert)
		}
	}
}
