/* 2019-01-01 (cc) <paul4hough@gmail.com>
   FIXME what is this for?
*/
package main

import (
	"fmt"
	"net/http"
)
type errHandler func(http.ResponseWriter, *http.Request) error

func (fn errHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if err := fn(w, r); err != nil {
		prom.Errors.Inc()
		fmt.Println("ERROR: ",err.Error())
    }
}
