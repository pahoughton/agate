/* 2019-01-01 (cc) <paul4hough@gmail.com>
   FIXME what is this for?
*/
package main

import (
	"fmt"
	"strings"
)

func procResolved(a *AmgrAlert) error {

	node := strings.Split(a.Labels["instance"],":")[0]

	tid, err := adb.GetTicket(a.StartsAt,node,a.Labels["alertname"])
	if err != nil {
		return err
	}
	com := fmt.Sprintf("resolved at %v",a.EndsAt)

	if err = addTicketComment(tid,com); err != nil {
		fmt.Println("ERROR: ticket comment - ",err.Error())
	}
	return adb.DelTicket(a.StartsAt,node,a.Labels["alertname"])
}
