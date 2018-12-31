/* 2018-12-26 (cc) <paul4hough@gmail.com>
   FIXME what is this for?
*/
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func handleShow(
	w http.ResponseWriter,
	r *http.Request ) {

	var tdata string

	tnumStr := r.URL.Query().Get("num")
	tnum, err := strconv.ParseUint(tnumStr, 10, 32)
	if err != nil {
		tdata = "<p><b>Error</b> num param(" +
			tnumStr + ") conv error: " + err.Error() + "</p>"
	} else if tnum >= uint64(len(tickets)) {
		tdata = fmt.Sprintf(
			"<p><b>Error</b> num out of range %d >= %d</p>",
			tnum, len(tickets))
	} else {
		fmt.Println(tickets[tnum])
		var tckMap map[string]string
		if err := json.Unmarshal([]byte(tickets[tnum]), &tckMap); err != nil {
			tdata = "<p><b>Error</b> json parse: " + err.Error() + "</p>" +
				" <p>" + tickets[tnum] + "</p>\n"
		} else {
			tdata = "<table>\n"
			for k, v := range  tckMap {
				tdata += fmt.Sprintf("<tr><td>%s:</td>" +
					"<td><pre>%s</pre></td><tr>\n",k,v)
			}
			tdata += "</table>\n"
		}
	}
	resp := fmt.Sprintf(`<!DOCTYPE html>
<html>
<body>
<h2> Ticket %d details </h2>

%s

</body>
</html>
`,
		tnum,
		tdata)

	w.WriteHeader(200)
	w.Write([]byte(resp))
}
