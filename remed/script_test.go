/* 2019-02-22 (cc) <paul4hough@gmail.com>
   FIXME what is this for?
*/
package remed

import (
	"io/ioutil"
	"os"
	"testing"

	// pmod "github.com/prometheus/common/model"
	"github.com/stretchr/testify/assert"
	"github.com/pahoughton/agate/config"
	"github.com/pahoughton/agate/amgr/alert"

)
func TestScriptAvail(t *testing.T) {
	cfg := config.New()
	cfg.Remed.ScriptsDir = "testdata/scripts"

	r := New(cfg.Remed,nil,false)
	labels := alert.LabelSet{}
	assert.False(t,r.ScriptAvail(labels))
	labels = alert.LabelSet{"alertname": "remed"}
	assert.True(t,r.ScriptAvail(labels))
	labels = alert.LabelSet{"alertname": "invalid"}
	assert.False(t,r.ScriptAvail(labels))

	r.Del()
}

func TestScript(t *testing.T) {
	cfg := config.New()
	cfg.Remed.ScriptsDir = "testdata/scripts"
	r := New(cfg.Remed,nil,false)
	assert.NotNil(t,r)

	tfn := "/tmp/test-agate-ansible"
	os.Remove(string(tfn))
	labels := alert.LabelSet{"alertname": "remed","testfn": tfn}
	got, err := r.Script("localhost", labels)
	if r.debug { print(got); }
	if r.debug && err != nil { print(err.Error()); }
	assert.Nil(t,err)
	assert.NotNil(t,got)
	//print("\n"+got+"\n")
	assert.FileExists(t,string(tfn))
	buf,err := ioutil.ReadFile(string(tfn))
	assert.Nil(t,err)
	assert.NotNil(t,buf)
	assert.Equal(t,`fixed`+"\n",string(buf))
	os.Remove(string(tfn))

	labels = alert.LabelSet{}
	got, err = r.Script("localhost", labels)
	assert.NotNil(t,err)

	labels = alert.LabelSet{"alertname": "invalid-mode","testfn": tfn}
	got, err = r.Script("localhost", labels)
	assert.NotNil(t,err)

	labels = alert.LabelSet{"alertname": "invalid","testfn": tfn}
	got, err = r.Script("localhost", labels)
	assert.NotNil(t,err)

	r.Del()
}
