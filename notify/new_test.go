/* 2019-03-31 (cc) <paul4hough@gmail.com>
   validate methods in new.go
*/
package notify

import (
	"os"
	"path/filepath"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/pahoughton/agate/config"
)

// clean data
func cleanNew(cfg config.Notify, dataDir string, dbg bool) *Notify {

	if files, err := filepath.Glob(filepath.Join(dataDir, "*-queue.bolt")); err == nil {
		for _, file := range files {
			err = os.RemoveAll(file)
			if err != nil {
				panic( err )
			}
		}
	} else {
		panic( err )
	}
	return New(cfg,dataDir,dbg)
}

func TestNew(t *testing.T) {
	cfg := config.New()
	got := New(cfg.Notify,"testdata",false)
	assert.NotNil(t,got)
	got.Del()
}

func TestNewPanic(t *testing.T) {
	assert.Panics(t, func() {
		cfg := config.New()
		cfg.Notify.Default = "george"
		New(cfg.Notify,"testdata",false)
	}, "New did not panic")
}

func TestValidSys(t *testing.T) {

	obj := New(config.New().Notify,"testdata",false)

	for _, sys := range []string{"mock","hpsm","gitlab"} {
		assert.True(t,obj.ValidSys(sys))
	}
	assert.False(t,obj.ValidSys("george"))
	obj.Del()
}

func TestSys(t *testing.T) {

	obj := New(config.New().Notify,"testdata",false)

	for _, sys := range []string{"mock","hpsm","gitlab"} {
		assert.NotNil(t,obj.Sys(sys))
	}

	assert.Nil(t,obj.Sys("george"))
	obj.Del()
}

func TestGroup(t *testing.T) {
	cfg := config.New()
	cfg.Notify.Sys.Gitlab.Group = "maul-alerts"
	cfg.Notify.Sys.Hpsm.Group = "WG1234"

	obj := New(cfg.Notify,"testdata",false)

	for _, sys := range []string{"hpsm","gitlab"} {
		assert.True(t,len(obj.Group(sys)) > 0 )
	}
	assert.True(t,len(obj.Group("george")) == 0 )
	obj.Del()
}
