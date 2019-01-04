/* 2018-12-21 (cc) <paul4hough@gmail.com>
   handle unknown http requests
*/
package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net/http"

	"github.com/boltdb/bolt"
)

func handleList(
	w http.ResponseWriter,
	r *http.Request ) {

	tckTable := "<table>\n"

	err := tdb.db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(Bucket)) // fixme skv bucket name

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {

			var t Ticket

			key := string(k)
			title := ""

			err := gob.NewDecoder(bytes.NewReader(v)).Decode(&t)
			if err != nil {
				fmt.Println("ERROR: ticket decode - ",err.Error())
				title = err.Error()
			} else {
				title = t.Title
			}

			tckTable += fmt.Sprintf(
			"<tr><td>%s</td>" +
				"<td><a href=\"http:/show?num=%s\">%s</a></td></tr>\n",
				key,key,title)
		}
		return nil
	})
	if err != nil {
		fmt.Println("ERROR: db view '",err.Error())
		tckTable += "<tr><td>ERROR: "+err.Error()+"</td></tr>\n"
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
