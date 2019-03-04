/* 2019-02-19 (cc) <paul4hough@gmail.com>
   FIXME what is this for?
*/
package db

import (
	"fmt"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestAGroupAdd(t *testing.T) {
	testit(t,func (t *testing.T,db *DB) {
		expStr := `{"json":"noise"}"`
		db.AGroupAdd([]byte(expStr),false)
	})
}
func TestAGroupQueue(t *testing.T) {
	testit(t,func (t *testing.T,db *DB) {
		expQsize := 16
		for qn := 0; qn < expQsize; qn += 1 {
			json := fmt.Sprintf(`{"qn":%d}`,qn)
			db.AGroupAdd([]byte(json),qn % (expQsize / 3) == 0)
		}
		agq := db.AGroupQueue()
		assert.Equal(t,expQsize,len(agq))
		for i, k := range agq {
			exp := fmt.Sprintf(`{"qn":%d}`,i)
			assert.Equal(t,exp,string(db.AGroupGet(k).Json))
			assert.Equal(t,i % (expQsize / 3) == 0,db.AGroupGet(k).Resolve)
		}
	})
}
func TestAGroupDel(t *testing.T) {
	testit(t,func (t *testing.T,db *DB) {
		expQsize := 16
		for qn := 0; qn < expQsize; qn += 1 {
			json := fmt.Sprintf(`{"qn":%d}`,qn)
			db.AGroupAdd([]byte(json),qn % (expQsize / 3) == 0)
		}
		agq := db.AGroupQueue()
		assert.Equal(t,expQsize,len(agq))
		for i, k := range agq {
			db.AGroupDel(k)
			assert.Equal(t,expQsize - i,len(db.AGroupQueue())+1)
		}
		assert.Equal(t,0,len(db.AGroupQueue()))
	})
}
