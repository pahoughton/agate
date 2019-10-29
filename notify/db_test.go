/* 2019-10-22 (cc) <paul4hough@gmail.com>
   FIXME what is this for?
*/
package notify

import (
	"testing"

	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/require"

	"github.com/pahoughton/agate/config"
	"github.com/pahoughton/agate/notify/note"
)
func TestDB(t *testing.T) {

	obj := New(config.New().Notify,"testdata",false)

	for _, sys := range []string{"mock","hpsm","gitlab"} {

		for _, grp := range []string{"alerts","warnings"} {

			assert.NotNil(t,obj.DB(sys,grp))
		}
	}
	akey := []byte("test-001")
	key := Key{"gitlab","alert",akey}
	note := note.Note{From: "test-01"}
	obj.dbUpdate(key,note)
	got := obj.dbGet(key)
	assert.Equal(t,note.From,got.From)
	got = obj.dbGet(Key{"other","alert",akey})
	assert.True(t,len(got.From) == 0)
	obj.dbDelete(key)
	got = obj.dbGet(key)
	assert.True(t,len(got.From) == 0)

	obj.Del()
}
