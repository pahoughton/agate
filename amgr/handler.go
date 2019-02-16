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
	am.metrics.AlertGroupsRecvd.With(
		promp.Labels{
			"resolve": resStr,
		}).Inc()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		am.Error(fmt.Sprintf(
			"amgr.ServeHTTP: ioutil.ReadAll - %s",err.Error()))
		return
	}
	defer r.Body.Close()
	err := am.db.AgroupAdd(b,resolve)
	if err != nil {
		am.Error(fmt.Sprintf(
			"amgr.ServeHTTP: db.AmgrAdd: %s",err.Error()))
		return
	}

	select {
	case a.manager <- true:
	case <- time.After(1):
	}
}
