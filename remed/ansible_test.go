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
func TestAnsibleAvail(t *testing.T) {
	cfg := config.New()
	cfg.Global.PlaybookDir = "testdata/playbook"

	r := New(cfg.Global,false)
	labels := pmod.LabelSet{}
	assert.False(t,r.AnsibleAvail(labels))
	labels = pmod.LabelSet{"alertname": "remed"}
	assert.True(t,r.AnsibleAvail(labels))
	labels = pmod.LabelSet{"alertname": "invalid"}
	assert.False(t,r.AnsibleAvail(labels))

	r.Close()
}

func TestAnsible(t *testing.T) {
	cfg := config.New()
	cfg.Global.PlaybookDir = "testdata/playbook"

	obj := New(cfg.Global,false)
	assert.NotNil(t,obj)

	tfn := pmod.LabelValue("/tmp/test-agate-ansible")
	os.Remove(string(tfn))
	labels := pmod.LabelSet{"alertname": "remed","testfn": tfn}
	got, err := obj.Ansible("localhost", labels)
	assert.Nil(t,err)
	if err != nil && len(got) > 0 && obj.debug { print("\n"+got+"\n") }
	assert.Regexp(t,"PLAY",got)

	assert.FileExists(t,string(tfn))
	buf,err := ioutil.ReadFile(string(tfn))
	assert.Nil(t,err)
	assert.NotNil(t,buf)
	assert.Equal(t,`"test"`+"\n",string(buf))
	os.Remove(string(tfn))

	labels = pmod.LabelSet{}
	got, err = obj.Ansible("localhost", labels)
	assert.NotNil(t,err)

	labels = pmod.LabelSet{"alertname": "invalid"}
	got, err = obj.Ansible("localhost", labels)
	assert.NotNil(t,err)

	obj.Close()
}
