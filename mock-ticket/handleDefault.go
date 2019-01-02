/* 2018-12-21 (cc) <paul4hough@gmail.com>
   handle unknown http requests
*/
package main

import (
	//	"bytes"
	//	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func handleDefault(
	w http.ResponseWriter,
	r *http.Request ) {

	// inc metrics counter
	prom.UnsupRecvd.Inc()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

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
		b)

	log.Warning("unsupported request\n",resp)
	w.WriteHeader(404)
	w.Write([]byte(resp))
}
