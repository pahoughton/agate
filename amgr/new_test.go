/* 2019-02-19 (cc) <paul4hough@gmail.com>
   FIXME what is this for?
*/
package amgr

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/pahoughton/agate/config"

)

func TestNew(t *testing.T) {
	got := New(config.New(),"testdata/data",false)
	assert.NotNil(t,got)
	got.Del()
	got = New(config.New(),"testdata/data",false)
	assert.NotNil(t,got)
	got.Del()
}
