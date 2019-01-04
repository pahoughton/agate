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
	"strings"

	promp "github.com/prometheus/client_golang/prometheus"
)

func handleAlertGroup(
	w http.ResponseWriter,
	r *http.Request ) error {

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("ioutil.ReadAll - %s",err.Error())
	}
	defer r.Body.Close()

	if *args.Debug {
		var dbgbuf bytes.Buffer
		if err := json.Indent(&dbgbuf, b, " ", "  "); err != nil {
			return fmt.Errorf("json.Indent: ",err.Error())
		}
		fmt.Println("DEBUG req body\n",dbgbuf.String())
	}

	var abody AmgrAlertBody
	if err := json.Unmarshal(b, &abody); err != nil {
		return fmt.Errorf("json.Unmarshal alert: %s\n%v",err.Error(),b)
    }

	prom.AlertGroupsRecvd.With(
		promp.Labels{
			"status": abody.Status,
			"receiver": abody.Receiver,
		}).Inc()

	for _, alert := range abody.Alerts {
		node := strings.Split(alert.Labels["instance"],":")[0]

		prom.AlertsRecvd.With(
			promp.Labels{
				"name": alert.Labels["alertname"],
				"node": node,
				"status": abody.Status,
			}).Inc()

		if alert.Status == "firing" {

			tid, err := createTicket(&alert);
			if err != nil {
				return fmt.Errorf("createTicket: %s",err.Error())
			}
			if _, ok := alert.Labels["ansible"]; ok {
				if err := procAnsible(&alert,tid); err != nil {
					return err
				}
			}
			if _, ok := alert.Labels["script"]; ok {
				if err := procScript(&alert,tid); err != nil {
					return err
				}
			}
		} else if alert.Status == "resolved" {
			return procResolved(&alert)
		} else {
			return fmt.Errorf("status: %s\n%v\n",alert.Status,alert)
		}
	}
	return nil
}
