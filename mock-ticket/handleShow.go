/* 2018-12-26 (cc) <paul4hough@gmail.com>
   FIXME what is this for?
*/
package main

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/boltdb/bolt"
)

func handleShow(
	w http.ResponseWriter,
	r *http.Request ) {

	var tdata string

	tid := r.URL.Query().Get("num")

	var t Ticket

	err := tdb.db.View(func(tx *bolt.Tx) error {
		bckt := tx.Bucket([]byte(Bucket))
		tgob := bckt.Get([]byte(tid))
		if tgob == nil {
			return errors.New("not found - "+tid)
		}
		return gob.NewDecoder(bytes.NewReader(tgob)).Decode(&t)
	})
	if err != nil {
		fmt.Println("ERROR: db GET '",err.Error())
		tdata = "<p><b>Error</b> get(" + tid + ") "+ err.Error() + "</p>"
	} else {
		tdata = `
<table>
  <tr><td>id</td><td>`+tid+`</td></tr>
  <tr><td>title</td><td>`+t.Title+`</td></tr>
  <tr><td>state</td><td>`+t.State+`</td></tr>
  <tr><td>node</td><td>`+t.Node+`</td></tr>
  <tr><td>worker</td><td>`+t.Worker+`</td></tr>
  <tr><td>desc</td><td><pre>`+t.Desc+`</pre></td></tr>`

		for i, cmt := range t.Comments {
			tdata += `
  <tr><td>comment `+strconv.Itoa(i+1)+"<td><pre>"+cmt+"</pre></td></tr>"
		}
		tdata += "\n</table>\n"
	}

	resp := fmt.Sprintf(`<!DOCTYPE html>
<html>
<body>
<h2> Ticket %s details </h2>

%s

</body>
</html>
`,
		tid,
		tdata)

	w.WriteHeader(200)
	w.Write([]byte(resp))
}
