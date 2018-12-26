/* 2018-12-25 (cc) <paul4hough@gmail.com>
   FIXME what is this for?
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

func handleUnsup(
	w http.ResponseWriter,
	r *http.Request ) {

	unsupRecvd.Inc()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()

	var buf bytes.Buffer
	if err := json.Indent(&buf, b, " >", "  "); err != nil {
		log.Fatal(err)
	}
	resp := fmt.Sprintf(`<!DOCTYPE html>
<html>
<body>
<h2> 404 Unsupported request </h2>

<b>remote:</b> %s <br/>
<b>host:</b> %s <br/>
<b>uri:</b> %s <br/>
<p><b>method:</b> %s<br/>
<p><b> payload </b>
<pre>
%s
</pre>

</body>
</html>
`,
		r.RemoteAddr,
		r.Host,
		r.RequestURI,
		r.Method,
		buf.String())
	log.Warning("unsupported request")
	fmt.Fprintf(os.Stderr,"req body\n%s\n",buf.String())
	w.WriteHeader(404)
	w.Write([]byte(resp))
}
