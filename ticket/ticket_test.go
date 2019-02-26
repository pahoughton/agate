/* 2019-02-17 (cc) <paul4hough@gmail.com>
   FIXME what is this for?
*/
package ticket

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/pahoughton/agate/config"
)

func TestNew(t *testing.T) {
	cfg := config.New()
	got := New(cfg.Ticket,false)
	assert.NotNil(t,got)
	got.Del()
}

func TestNewPanic(t *testing.T) {
	assert.Panics(t, func() {
		cfg := config.New()
		cfg.Ticket.Default = "george"
		New(cfg.Ticket,false)
	}, "New did not panic")
}
func TestSink(t *testing.T) {
	cfg := config.New()
	obj := New(cfg.Ticket,false)
	assert.NotNil(t,t)
	got := obj.Sink(TSysMock)
	assert.NotNil(t,got)
	assert.NotNil(t,got.Group())
	got = obj.Sink(TSysUnknown)
	assert.Nil(t,got)

	obj.Del()
}

func TestGroup(t *testing.T) {
	cfg := config.New()
	obj := New(cfg.Ticket,false)
	assert.NotNil(t,t)
	got := obj.Group(TSysMock)
	assert.Equal(t,got,"")
	obj.Del()

	exp := "agate-test"
	cfg.Ticket.Sys.Gitlab.Group = exp
	cfg.Ticket.Sys.Hpsm.Group = exp + "hpsm"
	obj = New(cfg.Ticket,false)
	assert.NotNil(t,t)

	assert.Equal(t,exp,obj.Group(TSysGitlab))
	assert.Equal(t,exp + "hpsm",obj.Group(TSysHpsm))
	obj.Del()

}
