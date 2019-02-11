/* 2019-01-03 (cc) <paul4hough@gmail.com>
   simple dump service
*/
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"gopkg.in/alecthomas/kingpin.v2"
)

func handleAny(
	w http.ResponseWriter,
	r *http.Request ) {

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("ERROR: ioutil.ReadAll - %s",err.Error())
		return
	}
	defer r.Body.Close()
	fmt.Printf("URL: %v\n",r.URL)
	for k,v := range r.Header {
		fmt.Printf("HDR %s: %v\n",k,v)
	}
	fmt.Println(string(b))
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
