/* 2019-02-22 (cc) <paul4hough@gmail.com>
   FIXME what is this for?
*/
package amgr

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/pahoughton/agate/config"
)
func TestDetect(t *testing.T) {
	os.Remove("testdata/data/agate.bolt")
	os.Remove("testdata/data/agate-1.bolt")
	cfg := config.New()
	am := New(cfg,"testdata/data",false)
	assert.NotNil(t,am)
	agq := am.db.AGroupQueueList(uint(am.notify.DefSys))
	assert.Equal(t,0,len(agq))

    rr := httptest.NewRecorder()
	bodyf, err := os.Open("testdata/amgr/alert-group.json")
	if err != nil { t.Fatal(err) }

	req, err := http.NewRequest("GET", "/alerts", bodyf)
    if err != nil { t.Fatal(err) }

	am.ServeHTTP(rr,req)
	bodyf.Close()
	assert.Equal(t,rr.Code,http.StatusOK)

	agq = am.db.AGroupQueueList(uint(am.notify.DefSys))
	assert.Equal(t,1,len(agq))


	bodyf, err = os.Open("testdata/amgr/alert-group-v3.json")
	if err != nil { t.Fatal(err) }

	req, err = http.NewRequest("GET", "/alerts", bodyf)
    if err != nil { t.Fatal(err) }

	assert.Panics(t, func() {
		am.ServeHTTP(rr,req)
	}, "detect.ServeHTTP did not panic")
	bodyf.Close()
	assert.Equal(t,rr.Code,500)

	bodyf, err = os.Open("testdata/amgr/alert-group-v5.json")
	if err != nil { t.Fatal(err) }

	req, err = http.NewRequest("GET", "/alerts", bodyf)
    if err != nil { t.Fatal(err) }

	assert.Panics(t, func() {
		am.ServeHTTP(rr,req)
	}, "detect.ServeHTTP did not panic")
	bodyf.Close()
	assert.Equal(t,rr.Code,500)

	am.Del()
}
