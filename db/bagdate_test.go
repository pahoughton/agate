/* 2019-02-19 (cc) <paul4hough@gmail.com>
   FIXME what is this for?
*/
package db

import (
	"fmt"
	"strconv"
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
)

var tbag = []byte("bag")
var tnow = time.Now()

func TestBagDateAdd(t *testing.T) {
	testit(t,func (t *testing.T,db *DB) {
		key := "abc"
		val := "def"
		db.BagDateAdd(tnow,tbag,[]byte(key),[]byte(val))
		assert.Nil(t,db.db.Close())
	})
}

func TestBagDateGet(t *testing.T) {
	testit(t,func (t *testing.T,db *DB) {
		dataCnt := 24
		data := make(map[string]string)
		for i := 0; i < dataCnt; i += 1 {
			data[fmt.Sprintf("k-%03d",i)] = strconv.Itoa(i)
		}
		for k, v := range data {
			db.BagDateAdd(tnow,tbag,[]byte(k),[]byte(v))
		}
		for k, v := range data {
			assert.Equal(t,[]byte(v),db.BagDateGet(tnow,tbag,[]byte(k)))
		}
	})
}
func TestBagDateDel(t *testing.T) {
	testit(t,func (t *testing.T,db *DB) {
		dataCnt := 24
		data := make(map[string]string)
		for i := 0; i < dataCnt; i += 1 {
			data[fmt.Sprintf("k-%03d",i)] = strconv.Itoa(i)
		}
		for k, v := range data {
			db.BagDateAdd(tnow,tbag,[]byte(k),[]byte(v))
		}
		for k, v := range data {
			assert.Equal(t,[]byte(v),db.BagDateGet(tnow,tbag,[]byte(k)))
			db.BagDateDel(tnow,tbag,[]byte(k))
			assert.Nil(t,db.BagDateGet(tnow,tbag,[]byte(k)))
		}
	})
}
