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

	got := New(cfg.Remed,nil,false)
	assert.NotNil(t,got)
	got.Del()
	got = New(cfg.Remed,nil,false)
	assert.NotNil(t,got)
	got.Del()
}

func TestNewPanic(t *testing.T) {
	var got *Remed
	var gerr *Remed
	assert.Panics(t, func() {
		cfg := config.New()
		got = New(cfg.Remed,nil,false)
		assert.NotNil(t,got)
		gerr = New(cfg.Remed,nil,false)
	}, "New did not panic")
	if got != nil { got.Del() }
	if gerr != nil { gerr.Del() }
}
