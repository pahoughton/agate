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

	resStr := r.FormValue("resolve")
	resolve := false
	if len(resStr) > 0 {
		resStr = "true"
		resolve = true
	} else {
		resStr = "false"
	}
	am.metrics.groups.With(promp.Labels{"resolve":resStr}).Inc()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		if am.debug {
			panic(err)
		} else {
			am.Error(err)
		}
		return
	}
	am.db.AGroupAdd(b,resolve)
	am.qmgr.Notify(1)
}
