/* 2019-02-22 (cc) <paul4hough@gmail.com>
   FIXME what is this for?
*/
package remed

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/pahoughton/agate/amgr/alert"
	"github.com/pahoughton/agate/config"

)
func TestAnsibleAvail(t *testing.T) {
	cfg := config.New()
	cfg.Remed.PlaybookDir = "testdata/playbook"

	r := New(cfg.Remed,nil,false)
	labels := alert.LabelSet{}
	assert.False(t,r.AnsibleAvail(labels))
	labels =
		alert.LabelSet{"alertname": "remed"}
	assert.True(t,r.AnsibleAvail(labels))
	labels = alert.LabelSet{"alertname": "invalid"}
	assert.False(t,r.AnsibleAvail(labels))

	r.Del()
}

func TestAnsible(t *testing.T) {
	if len(os.Getenv("TRAVIS")) > 0 || len(os.Getenv("GITLAB_CI")) > 0 {
		// travis and gitlab can't ssh localhost for ansible test
		print("travis/gitlab ci skip ansible\n")
		return
	}
	cfg := config.New()
	cfg.Remed.PlaybookDir = "testdata/playbook"

	obj := New(cfg.Remed,nil,false)
	assert.NotNil(t,obj)

	tfn := "/tmp/test-agate-ansible"
	os.Remove(string(tfn))
	labels := alert.LabelSet{"alertname": "remed","testfn": tfn}
	got, err := obj.Ansible("localhost", labels)
	assert.Nil(t,err)
	if err != nil && len(got) > 0 { print("\nERRgot: "+got+"\n") }
	assert.Regexp(t,"failed=0",got)

	assert.FileExists(t,string(tfn))
	buf,err := ioutil.ReadFile(string(tfn))
	assert.Nil(t,err)
	assert.NotNil(t,buf)
	assert.Equal(t,`"test"`+"\n",string(buf))
	os.Remove(string(tfn))

	labels = alert.LabelSet{}
	got, err = obj.Ansible("localhost", labels)
	assert.NotNil(t,err)

	labels = alert.LabelSet{"alertname": "invalid"}
	got, err = obj.Ansible("localhost", labels)
	assert.NotNil(t,err)

	obj.Del()
}
