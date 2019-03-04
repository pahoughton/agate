/* 2019-02-22 (cc) <paul4hough@gmail.com>
   FIXME what is this for?
*/
package remed

import (
	"io/ioutil"
	"os"
	"testing"

	pmod "github.com/prometheus/common/model"
	"github.com/stretchr/testify/assert"
	"github.com/pahoughton/agate/config"

)
func TestScriptAvail(t *testing.T) {
	cfg := config.New()
	cfg.Global.ScriptsDir = "testdata/scripts"

	r := New(cfg.Global,false)
	labels := pmod.LabelSet{}
	assert.False(t,r.ScriptAvail(labels))
	labels = pmod.LabelSet{"alertname": "remed"}
	assert.True(t,r.ScriptAvail(labels))
	labels = pmod.LabelSet{"alertname": "invalid"}
	assert.False(t,r.ScriptAvail(labels))

	r.Close()
}

func TestScript(t *testing.T) {
	cfg := config.New()
	cfg.Global.ScriptsDir = "testdata/scripts"
	r := New(cfg.Global,false)
	assert.NotNil(t,r)

	tfn := pmod.LabelValue("/tmp/test-agate-ansible")
	os.Remove(string(tfn))
	labels := pmod.LabelSet{"alertname": "remed","testfn": tfn}
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

	labels = pmod.LabelSet{}
	got, err = r.Script("localhost", labels)
	assert.NotNil(t,err)

	labels = pmod.LabelSet{"alertname": "invalid-mode","testfn": tfn}
	got, err = r.Script("localhost", labels)
	assert.NotNil(t,err)

	labels = pmod.LabelSet{"alertname": "invalid","testfn": tfn}
	got, err = r.Script("localhost", labels)
	assert.NotNil(t,err)

	r.Close()
}
