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
type CommandArgs struct {
	ListenAddr	*string
}

var (

	app = kingpin.New(filepath.Base(os.Args[0]),
		"mock dump http service").
			Version("0.1.1")

	args = CommandArgs{
		ListenAddr:	app.Flag("listen-addr","listen address").
			Default(":5101").String(),
	}
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
	fmt.Println(string(b))
}
func main() {
	http.HandleFunc("/",handleAny)

	fmt.Println("FATAL: ",http.ListenAndServe(*args.ListenAddr,nil).Error())
	os.Exit(1)
}
