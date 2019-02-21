/* 2019-02-19 (cc) <paul4hough@gmail.com>
   FIXME what is this for?
*/
package db

import (
	"errors"
	"fmt"
	"strconv"
	"testing"
	"time"
	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/assert"
)


func TestAlertAdd(t *testing.T) {
	testit(t,func (t *testing.T,db *DB) {
		key := "abc"
		val := "def"
		db.AlertAdd(time.Now(),[]byte(key),[]byte(val))
		assert.Nil(t,db.db.Close())
	})
}
func TestAlertGet(t *testing.T) {
	testit(t,func (t *testing.T,db *DB) {
		dataCnt := 24
		data := make(map[string]string)
		for i := 0; i < dataCnt; i += 1 {
			data[fmt.Sprintf("k-%03d",i)] = strconv.Itoa(i)
		}
		today := time.Now()
		for k, v := range data {
			db.AlertAdd(today,[]byte(k),[]byte(v))
		}
		for k, v := range data {
			assert.Equal(t,[]byte(v),db.AlertGet(today,[]byte(k)))
		}
	})
}
func TestAlertDel(t *testing.T) {
	testit(t,func (t *testing.T,db *DB) {
		dataCnt := 24
		data := make(map[string]string)
		for i := 0; i < dataCnt; i += 1 {
			data[fmt.Sprintf("k-%03d",i)] = strconv.Itoa(i)
		}
		today := time.Now()
		for k, v := range data {
			db.AlertAdd(today,[]byte(k),[]byte(v))
		}
		for k, v := range data {
			assert.Equal(t,[]byte(v),db.AlertGet(today,[]byte(k)))
			db.AlertDel(today,[]byte(k))
			assert.Nil(t,db.AlertGet(today,[]byte(k)))
		}
	})
}
func TestAlertCleanBuckets(t *testing.T) {
	testit(t,func (t *testing.T,db *DB) {
		maxDays := db.maxDays
		today := time.Now()
		for i := 0; i < maxDays * 2; i += 1 {
			db.AlertAdd(today.AddDate(0,0,i * -1),[]byte("dont"),[]byte("care"))
		}
		db.AlertCleanBuckets()
		for i := 0; i < maxDays; i += 1 {

			bname := today.AddDate(0,0,i * -1).Format(BNameFmt)
			err := db.db.View(func(tx *bolt.Tx) error {
				if ab := tx.Bucket(alertsBucket()); ab != nil {
					if b := ab.Bucket([]byte(bname)); b != nil {
						return nil
					}
				}
				return errors.New(bname)
			})
			assert.Nil(t,err)
		}
		for i := int(maxDays+1); i < int(maxDays * 2) + 7; i += 1 {

			bname := today.AddDate(0,0,i * -1).Format(BNameFmt)
			err := db.db.View(func(tx *bolt.Tx) error {
				if ab := tx.Bucket(alertsBucket()); ab != nil {
					if b := ab.Bucket([]byte(bname)); b == nil {
						return nil
					}
				}
				return errors.New(bname)
			})
			assert.Nil(t,err)
		}
	})
}
