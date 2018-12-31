/* 2018-12-21 (cc) <paul4hough@gmail.com>
   handle unknown http requests
*/
package main

import (
	//	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func handleList(
	w http.ResponseWriter,
	r *http.Request ) {

	tckTable := "<table>\n"
	for i,tck := range tickets {
		var tckMap map[string]string
		if err := json.Unmarshal([]byte(tck), &tckMap); err != nil {
			log.Error(err)
			return
		}
		tckTable += fmt.Sprintf(
			"<tr><td>%d</td>" +
				"<td><a href=\"http:/show?num=%d\">%s</a></td></tr>\n",
			i, i, tckMap["title"])
	}
	tckTable += "</table>\n"

	resp := fmt.Sprintf(`<!DOCTYPE html>
<html>
<body>
<h2> Ticket List </h2>

%s

</body>
</html>
`,
		tckTable)

	fmt.Println(resp)
	w.WriteHeader(200)
	w.Write([]byte(resp))
}
