/* 2019-02-22 (cc) <paul4hough@gmail.com>
   FIXME what is this for?
*/
package amgr

import (
	"net/http/httptest"
	"os"
	"testing"

	pmod "github.com/prometheus/common/model"
	"github.com/stretchr/testify/assert"
	"github.com/pahoughton/agate/amgr/alert"
	"github.com/pahoughton/agate/config"
	"github.com/pahoughton/agate/ticket/mock"

)
func TestRemed(t *testing.T) {
	tsmock := &mock.MockServer{}
	svcmock := httptest.NewServer(tsmock)
	defer svcmock.Close()

	cfg := config.New()
	cfg.Global.ScriptsDir = "testdata/scripts"
	cfg.Global.PlaybookDir = "testdata/playbook"
	cfg.Ticket.Sys.Mock.Url = svcmock.URL
	am := New(cfg,"testdata/data",false)

	assert.NotNil(t,am)

	tid := am.ticket.Create(am.ticket.Default,"grp","title","desc")

	expfn := "/tmp/test-agate-ansible"
	os.Remove(expfn)
	alert := alert.Alert{}
	alert.Labels = pmod.LabelSet{
		"alertname": "remed",
		"instance": "localhost:9100",
		"testfn": pmod.LabelValue(expfn),
	}
	exphits := tsmock.Hits + 1
	if len(os.Getenv("TRAVIS")) == 0 {
		// travis can't ssh localhost for ansible test
		am.Remed(alert,tid)
		am.fix.wg.Wait()
		assert.Equal(t,exphits,tsmock.Hits)
		assert.FileExists(t,expfn)
		os.Remove(expfn)
	} else {
		print("travis - skip ansible")
	}

	expfn = "/tmp/test-agate-script"
	os.Remove(expfn)
	alert.Labels = pmod.LabelSet{
		"alertname": "fix",
		"instance": "localhost:9100",
		"testfn": pmod.LabelValue(expfn),
	}

	exphits = tsmock.Hits + 1
	am.Remed(alert,tid)
	am.fix.wg.Wait()
	assert.Equal(t,exphits,tsmock.Hits)
	assert.FileExists(t,expfn)
	os.Remove(expfn)

	alert.Labels = pmod.LabelSet{
		"alertname": "invalid",
		"instance": "localhost:9100",
	}

	exphits = tsmock.Hits
	am.Remed(alert,tid)
	am.fix.wg.Wait()
	assert.Equal(t,exphits,tsmock.Hits)

	am.Close()
}
