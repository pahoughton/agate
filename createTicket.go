/* 2018-12-27 (cc) <paul4hough@gmail.com>
   FIXME what is this for?
*/
package main

import "errors"

func createTicket(a *AmgrAlert) (string, error) {

	if args.TicketURL != nil {
		return createApiTicket(a)
	} else if args.SMTPAddr != nil &&
			args.EmailTo != nil &&
			args.EmailFrom != nil {
		return createEmailTicket(a)
	} else {
		return "", errors.New("no ticket destination")
	}
}
