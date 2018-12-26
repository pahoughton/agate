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

	log "github.com/sirupsen/logrus"
)

func handleAlertGroup(
	w http.ResponseWriter,
	r *http.Request ) {

	alertGroupsRecvd.Inc()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()

	// fixme debug only
	var dbgbuf bytes.Buffer
	if err := json.Indent(&dbgbuf, b, " >", "  "); err != nil {
		log.Fatal(err)
	}
	log.Debug("handlerAlertGroup")
	fmt.Fprintf(os.Stderr,"DEBUG req body\n%s\n",dbgbuf.String())

	var abody AmgrAlertBody
	if err := json.Unmarshal(b, &abody); err != nil {
        log.Fatal(err)
    }

	if abody.Version != "4" {
		log.Fatal("unsupported json version: " + abody.Version)
	}
	// ignore resolved
	if abody.Status == "resolved" {
		resolvedGroupsRecvd.Inc()
		return
	}
	if abody.Status != "firing" {
		log.Fatal("unexpeded alert status: " + abody.Status)
	}
	for _, alert := range abody.Alerts {
		alertsRecvd.Inc()
		createTicket(&alert)
		if _, ok := alert.Labels["script"]; ok {
			procScript(&alert)
		}
		if _, ok := alert.Labels["ansible"]; ok {
			procAnsible(&alert)
		}
	}
}
