/* 2018-12-21 (cc) <paul4hough@gmail.com>
   handle ticket http requests
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

func handleTicket(
	w http.ResponseWriter,
	r *http.Request ) {

	// inc metrics counter
	ticketRecvd.Inc()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	var dbgbuf bytes.Buffer
	if err := json.Indent(&dbgbuf, b, " >", "  "); err != nil {
		log.Fatal(err)
	}

	tickets = append(tickets,string(b))

	resp := fmt.Sprintf(`<!DOCTYPE html>
<html>
<body>
<h2> Ticket Created </h2>
<p><b> payload </b></p>
<pre>
%s
</pre>
</body>
</html>
`,b)
	log.Info("ticket created")
	log.Debug("ticket resp")
	fmt.Fprintf(os.Stderr,"DEBUG req body\n%s\n",dbgbuf.String())
	w.WriteHeader(200)
	w.Write([]byte(resp))
}
