/* 2018-12-25 (cc) <paul4hough@gmail.com>
   create ticket from alert
*/
package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

func createTicket(a *AmgrAlert){

	log.Debug("generating ticket for alert: ")
	fmt.Fprintf(os.Stderr,"%+v\n",*a)

	ticketGend.Inc()
}
