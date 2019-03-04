/* 2019-02-14 (cc) <paul4hough@gmail.com>
   alertmanager alert handler
   add body to data store and notify manager
*/
package amgr

import (
	"io/ioutil"
	"net/http"
	promp "github.com/prometheus/client_golang/prometheus"
)

func (am *Amgr)ServeHTTP(w http.ResponseWriter,r *http.Request) {

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		if am.debug {
			panic(err)
		} else {
			am.Error(err)
		}
		w.WriteHeader(500)
		return
	}
	resStr := r.FormValue("resolve")
	resolve := false
	if len(resStr) > 0 {
		resStr = "true"
		resolve = true
	} else {
		resStr = "false"
	}
	am.metrics.groups.With(promp.Labels{"resolve":resStr}).Inc()

	if len(b) > 0 {
		am.db.AGroupAdd(b,resolve)
		am.qmgr.Notify(1)
		w.WriteHeader(200)
	} else {
		w.WriteHeader(500)
	}
}
