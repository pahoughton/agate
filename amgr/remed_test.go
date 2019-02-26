/* 2019-02-22 (cc) <paul4hough@gmail.com>
   FIXME what is this for?
*/
package amgr

import (
	"os"
	"testing"

	pmod "github.com/prometheus/common/model"
	"github.com/stretchr/testify/assert"
	"github.com/pahoughton/agate/amgr/alert"
	"github.com/pahoughton/agate/config"

)
func TestFix(t *testing.T) {
	cfg := config.New()
	cfg.Global.ScriptsDir = "testdata/scripts"
	cfg.Global.PlaybookDir = "testdata/playbook"
	am := New(cfg,"testdata/data",false)
	assert.NotNil(t,am)

	alert := alert.Alert{}
	atfn := pmod.LabelValue("/tmp/test-agate-ansible")
	os.Remove(string(atfn))
	alert.Labels = pmod.LabelSet{
		"alertname": "remed",
		"instance": "localhost:9100",
		"testfn": atfn,
	}

/*
	got := am.Fix(alert)
	assert.True(t,len(got) > 0)
	assert.FileExists(t,string(atfn))
	os.Remove(string(atfn))

	stfn := pmod.LabelValue("/tmp/test-agate-script")
	os.Remove(string(stfn))
	alert.Labels = pmod.LabelSet{
		"alertname": "fix",
		"instance": "localhost:9100",
		"testfn": stfn,
	}

	got = am.Fix(alert)
	assert.True(t,len(got) > 0)
	assert.FileExists(t,string(stfn))
	os.Remove(string(stfn))
	alert.Labels = pmod.LabelSet{
		"alertname": "invalid",
		"instance": "localhost:9100",
	}
	got = am.Fix(alert)
	assert.True(t,len(got) == 0)
*/
	am.Close()
}
