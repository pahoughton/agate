/* 2019-02-14 (cc) <paul4hough@gmail.com>
   alertmanager alert handler
   add body to data store and notify manager
*/
package amgr

func (am *Amgr)ServeHTTP(w http.ResponseWriter,r *http.Request) {

	resStr := r.FormValue("resolve")
	resolve := false
	if len(resStr) > 0 {
		resStr = "true"
		resolve = true
	} else {
		resStr = "false"
	}
	am.metrics.Recvd.With(promp.Labels{"resolve":resStr}).Inc()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		if am.debug {
			panic(err)
		} else {
			am.ErrorMesg(err,"amgr.ServeHTTP: ioutil.ReadAll")
		}
		return
	}
	if err := am.db.AgroupAdd(b,resolve); err != nil {
		panic(err)
	}
	am.manager.Notify(1)
}
