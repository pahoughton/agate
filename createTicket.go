/* 2018-12-27 (cc) <paul4hough@gmail.com>
   FIXME what is this for?
*/
package main

import "fmt"

func createTicket(a *AmgrAlert){

	if args.TicketMockURL != nil {
		createMockApiTicket(a)
	} else if args.SMTPAddr != nil &&
			args.EmailTo != nil &&
			args.EmailFrom != nil {
		createEmailTicket(a)
	} else {
		fmt.Println("FATAL: missing ticket destination")
	}
}
