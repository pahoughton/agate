/* 2019-02-22 (cc) <paul4hough@gmail.com>
   FIXME what is this for?
*/
package amgr

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/pahoughton/agate/db"
	"github.com/pahoughton/agate/config"
	"github.com/pahoughton/agate/notify"
	"github.com/pahoughton/agate/amgr/alert"
)

const tdir = "testdata/data"

func TestDetect(t *testing.T) {
	for _, pfn := range db.DbPrevFn {
		os.Remove(path.Join(tdir,pfn))
	}
	os.Remove(path.Join(tdir,db.DbFn))
	cfg := config.New()
	am := New(cfg,tdir,false)
	assert.NotNil(t,am)
	agq := am.db.AGroupQueueList(uint(am.notify.DefSys))
	assert.Equal(t,0,len(agq))

    rr := httptest.NewRecorder()
	bodyf, err := os.Open("testdata/amgr/alert-group.json")
	if err != nil { t.Fatal(err) }

	req, err := http.NewRequest("GET", Url, bodyf)
    if err != nil { t.Fatal(err) }

	am.ServeHTTP(rr,req)
	bodyf.Close()
	assert.Equal(t,rr.Code,http.StatusOK)

	agq = am.db.AGroupQueueList(uint(am.notify.DefSys))
	assert.Equal(t,1,len(agq))


	// unsupported version
	bodyf, err = os.Open("testdata/amgr/alert-group-v3.json")
	if err != nil { t.Fatal(err) }

	req, err = http.NewRequest("GET", Url, bodyf)
    if err != nil { t.Fatal(err) }

	assert.Panics(t, func() {
		am.ServeHTTP(rr,req)
	}, "detect.ServeHTTP did not panic")
	bodyf.Close()
	assert.Equal(t,rr.Code,500)

	bodyf, err = os.Open("testdata/amgr/alert-group-v5.json")
	if err != nil { t.Fatal(err) }

	req, err = http.NewRequest("GET", Url, bodyf)
    if err != nil { t.Fatal(err) }

	assert.Panics(t, func() {
		am.ServeHTTP(rr,req)
	}, "detect.ServeHTTP did not panic")
	bodyf.Close()
	assert.Equal(t,rr.Code,500)

	am.Del()
}

func TestDetectParams(t *testing.T) {
	for _, pfn := range db.DbPrevFn {
		os.Remove(path.Join(tdir,pfn))
	}
	os.Remove(path.Join(tdir,db.DbFn))
	cfg := config.New()
	am := New(cfg,tdir,false)
	assert.NotNil(t,am)
	agq := am.db.AGroupQueueList(uint(am.notify.DefSys))
	assert.Equal(t,0,len(agq))

    rr := httptest.NewRecorder()
	bodyf, err := os.Open("testdata/amgr/alert-group.json")
	if err != nil { t.Fatal(err) }

	// param ?system=gitlab
    nsysid := notify.NSysGitlab
	url :=  Url + "?system=gitlab"
	req, err := http.NewRequest("GET", url, bodyf)
    if err != nil { t.Fatal(err) }

	am.ServeHTTP(rr,req)
	bodyf.Close()
	assert.Equal(t,rr.Code,http.StatusOK)

	agq = am.db.AGroupQueueList(uint(nsysid))
	assert.Equal(t,1,len(agq))
	ag := alert.NewAlertGroup(am.db.AGroupQueueGet(uint(nsysid),agq[0]))
	require.NotNil(t,ag)
	nsys := am.db.AGroupNSysGet(ag.StartsAt(),ag.Key())
	require.NotNil(t,nsys)
	am.db.AGroupNSysDel(ag.StartsAt(),ag.Key())
	am.db.AGroupNidDel(ag.StartsAt(),ag.Key())
	am.db.AGroupQueueDel(uint(nsysid),agq[0])

	// param ?group=maul/alerts
	exp := "maul/alerts"
	url =  Url + "?system=gitlab&group=" + exp
	bodyf, err = os.Open("testdata/amgr/alert-group.json")
	req, err = http.NewRequest("GET",url,bodyf)
    if err != nil { t.Fatal(err) }

	am.ServeHTTP(rr,req)
	bodyf.Close()
	assert.Equal(t,rr.Code,http.StatusOK)

	agq = am.db.AGroupQueueList(uint(nsysid))
	assert.Equal(t,1,len(agq))
	ag = alert.NewAlertGroup(am.db.AGroupQueueGet(uint(nsysid),agq[0]))
	require.NotNil(t,ag)
	nsys = am.db.AGroupNSysGet(ag.StartsAt(),ag.Key())
	require.NotNil(t,nsys)
	assert.Equal(t,exp,nsys.Grp)
	am.db.AGroupNSysDel(ag.StartsAt(),ag.Key())
	am.db.AGroupNidDel(ag.StartsAt(),ag.Key())
	am.db.AGroupQueueDel(uint(nsysid),agq[0])

	// param ?no_resolve=true
    nsysid = am.notify.DefSys
	url =  Url + "?no_resolve=true"
	bodyf, err = os.Open("testdata/amgr/alert-group.json")
	req, err = http.NewRequest("GET",url,bodyf)
    if err != nil { t.Fatal(err) }

	am.ServeHTTP(rr,req)
	bodyf.Close()
	assert.Equal(t,rr.Code,http.StatusOK)

	agq = am.db.AGroupQueueList(uint(nsysid))
	assert.Equal(t,1,len(agq))
	ag = alert.NewAlertGroup(am.db.AGroupQueueGet(uint(nsysid),agq[0]))
	require.NotNil(t,ag)
	nsys = am.db.AGroupNSysGet(ag.StartsAt(),ag.Key())
	require.NotNil(t,nsys)
	assert.Equal(t,false,nsys.Resolve)

	am.Del()
}
