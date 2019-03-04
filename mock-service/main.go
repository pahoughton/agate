/* 2019-01-03 (cc) <paul4hough@gmail.com>
   simple dump service
*/
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	TimeFmt = "2006-01-02.15:04:05.000000000-07:00"
)
func handleAny(
	w http.ResponseWriter,
	r *http.Request ) {

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("ERROR: ioutil.ReadAll - %s",err.Error())
		return
	}
	output := ""

	defer r.Body.Close()
	fmt.Println(time.Now().Format(TimeFmt))
	fmt.Printf("URL: %v\n",r.URL)
	for k,v := range r.Header {
		output += fmt.Sprintf("HDR %s: %v\n",k,v)
	}
	jsonOut := false
	if atypeList, ok := r.Header["Content-Type"]; ok {
		for _, atype := range atypeList {
			if atype == "application/json" {
				var dbgbuf bytes.Buffer
				if err := json.Indent(&dbgbuf, b, "", "  "); err != nil {
					output += fmt.Sprintf("ERROR json.Indent: %s",err.Error())
				} else {
					output += dbgbuf.String()
					jsonOut = true
				}
				break
			}
		}
	}
	if ! jsonOut {
		output += string(b) + "\n"
	}
	fmt.Println(output)
}

func main() {
	app := kingpin.New(filepath.Base(os.Args[0]),
		"mock dump http service").
			Version("0.1.1")

	laddr := app.Flag("listen-addr","listen address").
		Default(":5101").String()

	kingpin.MustParse(app.Parse(os.Args[1:]))

	http.HandleFunc("/",handleAny)

	fmt.Println("FATAL: ",http.ListenAndServe(*laddr,nil).Error())
	os.Exit(1)
}
