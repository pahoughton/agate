/* 2019-02-17 (cc) <paul4hough@gmail.com>
   validate mocktid interface
*/
package tid

import (
	"testing"
	"github.com/stretchr/testify/assert"
)
func TestNewBytes(t *testing.T) {
	exp := []byte("secret-sauce")
	got := NewBytes(exp)
	assert.Equal(t,exp,got.Bytes())
	assert.Equal(t,string(exp),got.String())
}

func TestNewString(t *testing.T) {
	exp := "secret-sauce"
	got := NewString(exp)
	assert.Equal(t,exp,got.String())
	assert.Equal(t,[]byte(exp),got.Bytes())
}
