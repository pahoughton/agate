/* 2019-02-22 (cc) <paul4hough@gmail.com>
*/
package remed

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	// amgrntfy "github.com/prometheus/alertmanager/notify"
	amgrtmpl "github.com/prometheus/alertmanager/template"

	"github.com/pahoughton/agate/config"
	"github.com/pahoughton/agate/notify"
	"github.com/pahoughton/agate/notify/mock"
	"github.com/pahoughton/agate/amgr/alert"
)

func TestAlertHasRemed(t *testing.T) {
	cfg := config.New()
	cfg.Remed.PlaybookDir = "testdata/playbook"
	cfg.Remed.ScriptsDir = "testdata/scripts"

	r := New(cfg.Remed,nil,false)
	a := alert.Alert{}
	a.Labels = make(amgrtmpl.KV)
	assert.False(t,r.AlertHasRemed(a))
	a.Labels["alertname"] = "remed"
	assert.True(t,r.AlertHasRemed(a))
	a.Labels["alertname"] = "fix"
	assert.True(t,r.AlertHasRemed(a))
	a.Labels["alertname"] = "invalid"
	assert.False(t,r.AlertHasRemed(a))

	r.Del()
}

func TestAlertGroupHasRemed(t *testing.T) {
	cfg := config.New()
	cfg.Remed.PlaybookDir = "testdata/playbook"
	cfg.Remed.ScriptsDir = "testdata/scripts"

	r := New(cfg.Remed,nil,false)
	a := amgrtmpl.Alert{}
	a.Labels = make(amgrtmpl.KV)

	ag := alert.AlertGroup{	Data: &amgrtmpl.Data{ Alerts: []amgrtmpl.Alert{a}}}

	assert.False(t,r.AGroupHasRemed(ag))
	a.Labels["alertname"] = "remed"
	assert.True(t,r.AGroupHasRemed(ag))
	a.Labels["alertname"] = "fix"
	assert.True(t,r.AGroupHasRemed(ag))
	a.Labels["alertname"] = "invalid"
	assert.False(t,r.AGroupHasRemed(ag))

	r.Del()
}

func TestRemed(t *testing.T) {
	tsmock := &mock.MockServer{}
	svcmock := httptest.NewServer(tsmock)
	defer svcmock.Close()

	cfg := config.New()
	cfg.Notify.Sys.Mock.Url = svcmock.URL
	n := notify.New(cfg.Notify,false)
	assert.NotNil(t,n)

	cfg.Remed.ScriptsDir = "testdata/scripts"
	cfg.Remed.PlaybookDir = "testdata/playbook"
	r := New(cfg.Remed,n,false)
	assert.NotNil(t,r)

	nid, err := r.notify.Create(notify.NSysMock,"grp","title","desc",false,false)
	require.Nil(t,err)
	require.NotNil(t,nid)

	expfn := "/tmp/test-agate-ansible"
	os.Remove(expfn)
	alert := alert.Alert{}
	alert.Labels = amgrtmpl.KV{
		"alertname": "remed",
		"instance": "localhost:9100",
		"testfn": expfn,
	}

	exphits := tsmock.Hits + 1
	if len(os.Getenv("TRAVIS")) == 0 {
		// travis can't ssh localhost for ansible test
		r.Remed(alert,nid)
		r.wg.Wait()
		assert.Equal(t,exphits,tsmock.Hits)
		assert.FileExists(t,expfn)
		os.Remove(expfn)
	} else {
		print("travis - skip ansible")
	}

	expfn = "/tmp/test-agate-script"
	os.Remove(expfn)
	alert.Labels = amgrtmpl.KV{
		"alertname": "fix",
		"instance": "localhost:9100",
		"testfn": expfn,
	}

	exphits = tsmock.Hits + 1
	r.Remed(alert,nid)
	r.wg.Wait()
	assert.Equal(t,exphits,tsmock.Hits)
	assert.FileExists(t,expfn)
	os.Remove(expfn)

	alert.Labels = amgrtmpl.KV{
		"alertname": "invalid",
		"instance": "localhost:9100",
	}

	exphits = tsmock.Hits
	r.Remed(alert,nid)
	r.wg.Wait()
	assert.Equal(t,exphits,tsmock.Hits)

	r.Del()
}
