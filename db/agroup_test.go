/* 2019-02-19 (cc) <paul4hough@gmail.com>
   FIXME what is this for?
*/
package db

import (
	"fmt"
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var tnsys NSys = NSys{Sys: 1, Grp: "agate", Resolve: true}

func TestAGroupQueueNSysAdd(t *testing.T) {
	testit(t,func (t *testing.T,db *DB) {
		db.AGroupQueueNSysAdd(time.Now(),tnsys,[]byte("key"),[]byte("val"))
		assert.Equal(t,1,len(db.AGroupQueueList(tnsys.Sys)))
	})
}

func TestAGroupNSysGet(t *testing.T) {
	testit(t,func (t *testing.T,db *DB) {
		key := []byte("key")
		exp := `{"json":"noise"}"`
		db.AGroupQueueNSysAdd(tnow,tnsys,key,[]byte(exp))
		assert.Equal(t,1,len(db.AGroupQueueList(tnsys.Sys)))
		got := db.AGroupNSysGet(tnow,key)
		require.NotNil(t,got)
		assert.Equal(t,tnsys,*got)
	})
}
func TestAGroupNSysDel(t *testing.T) {
	testit(t,func (t *testing.T,db *DB) {
		key := []byte("key")
		exp := `{"json":"noise"}"`
		db.AGroupQueueNSysAdd(tnow,tnsys,key,[]byte(exp))
		got := db.AGroupNSysGet(tnow,key)
		assert.Equal(t,tnsys,*got)
		db.AGroupNSysDel(tnow,key)
		got = db.AGroupNSysGet(tnow,key)
		assert.Nil(t,got)
	})
}

func TestAGroupQueueAdd(t *testing.T) {
	testit(t,func (t *testing.T,db *DB) {
		db.AGroupQueueNSysAdd(time.Now(),tnsys,[]byte("key"),[]byte("val"))
	})
}

func TestAGroupQueueList(t *testing.T) {
	testit(t,func (t *testing.T,db *DB) {
		expQsize := 16
		for qn := 0; qn < expQsize; qn += 1 {
			val := fmt.Sprintf(`{"qn":%d}`,qn)
			db.AGroupQueueAdd(tnsys.Sys,[]byte(val))
		}
		agq := db.AGroupQueueList(tnsys.Sys)
		assert.Equal(t,expQsize,len(agq))
		for i, k := range agq {
			exp := fmt.Sprintf(`{"qn":%d}`,i)
			assert.Equal(t,exp,string(db.AGroupQueueGet(tnsys.Sys,k)))
		}
	})
}

func TestAGroupQueueDel(t *testing.T) {
	testit(t,func (t *testing.T,db *DB) {
		expQsize := 16
		for qn := 0; qn < expQsize; qn += 1 {
			val := fmt.Sprintf(`{"qn":%d}`,qn)
			db.AGroupQueueAdd(tnsys.Sys,[]byte(val))
		}
		agq := db.AGroupQueueList(tnsys.Sys)
		assert.Equal(t,expQsize,len(agq))
		for i, k := range agq {
			db.AGroupQueueDel(tnsys.Sys,k)
			assert.Equal(t,expQsize - i,len(db.AGroupQueueList(tnsys.Sys))+1)
		}
		assert.Equal(t,0,len(db.AGroupQueueList(tnsys.Sys)))
	})
}
