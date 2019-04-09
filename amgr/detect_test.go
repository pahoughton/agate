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

	am.Del()
}
