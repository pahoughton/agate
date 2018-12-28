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
)

func handleUnsup(
	w http.ResponseWriter,
	r *http.Request ) {

	prom.UnsupRecvd.Inc()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("FATAL-ioutil.ReadAll: %s",err.Error())
		os.Exit(2)
	}
	defer r.Body.Close()

	var buf bytes.Buffer
	if err := json.Indent(&buf, b, " >", "  "); err != nil {
		fmt.Println("FATAL-json.Indent: ",err.Error())
		os.Exit(2)
	}
	if *args.Debug {
		fmt.Printf("DEBUG req body\n%s\n",buf.String())
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
	fmt.Printf("WARN-unsupported\n%s\n",resp)
	w.WriteHeader(404)
	w.Write([]byte(resp))
}
