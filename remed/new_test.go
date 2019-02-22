/* 2019-02-19 (cc) <paul4hough@gmail.com>
   FIXME what is this for?
*/
package remed

import (
	"testing"

	"github.com/pahoughton/agate/config"
	"github.com/stretchr/testify/assert"
)
func TestNew(t *testing.T) {
	cfg := config.New()

	got := New(cfg.Global,false)
	assert.NotNil(t,got)
	got.Close()
	got = New(cfg.Global,false)
	assert.NotNil(t,got)
	got.Close()
}

func TestNewPanic(t *testing.T) {
	var got *Remed
	assert.Panics(t, func() {
		cfg := config.New()
		got = New(cfg.Global,false)
		New(cfg.Global,false)
	}, "New did not panic")
	if got != nil { got.Close() }
}
