/* 2019-02-19 (cc) <paul4hough@gmail.com>
   FIXME what is this for?
*/
package amgr

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	am := New(config.New(),false)
	assert.NotNil(am)
}
