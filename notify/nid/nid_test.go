/* 2019-02-17 (cc) <paul4hough@gmail.com>
   validate mocktid interface
*/
package nid

import (
	"testing"
	"github.com/stretchr/testify/assert"
)
func TestNewBytes(t *testing.T) {
	expsys := uint8(3)
	expstr := "secret-sauce"
	exp := append([]byte(expstr),byte(expsys))
	got := NewBytes(exp)
	assert.Equal(t,exp,got.Bytes())
	assert.Equal(t,expstr,got.Id())
	assert.Equal(t,uint(expsys),got.Sys())
}

func TestNewString(t *testing.T) {
	expsys := uint8(16)
	exp := "secret-sauce"
	got := NewString(expsys,exp)
	assert.Equal(t,exp,got.Id())
	assert.Equal(t,uint(expsys),got.Sys())
	assert.Equal(t,append([]byte(exp),byte(expsys)),got.Bytes())
}
