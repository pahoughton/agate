/* 2019-03-31 (cc) <paul4hough@gmail.com>
   validate methods in new.go
*/
package notify

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/pahoughton/agate/config"

)
type TestNSys struct {
	exp		NSys
	val		string
}

func TestNewNSys(t *testing.T) {
	tlist := []TestNSys{
		{ NSysMock, "mock" },
		{ NSysGitlab, "gitlab" },
		{ NSysHpsm, "hpsm"  },
		{ NSysUnknown, "invalid" },
	}
	for _, vt := range tlist {
		assert.Equal(t,vt.exp,NewNSys(vt.val))
	}
}
func TestNSysString(t *testing.T) {
	for exp, got := range nsysmap {
		assert.Equal(t,exp,got.String())
	}
	got := NSys(5)
	assert.Equal(t,"invalid",got.String())
}

func TestNew(t *testing.T) {
	cfg := config.New()
	got := New(cfg.Notify,false)
	assert.NotNil(t,got)
	got.Del()
}

func TestNewPanic(t *testing.T) {
	assert.Panics(t, func() {
		cfg := config.New()
		cfg.Notify.Default = "george"
		New(cfg.Notify,false)
	}, "New did not panic")
}
func TestSystem(t *testing.T) {
	cfg := config.New()
	obj := New(cfg.Notify,false)
	assert.NotNil(t,obj)
	for i := NSysMock; i < NSysUnknown; i += 1 {
		assert.NotNil(t,obj.System(i))
	}
	assert.Nil(t,obj.System(124))
	obj.Del()
}
