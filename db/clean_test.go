/* 2019-04-07 (cc) <paul4hough@gmail.com>
*/
package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/boltdb/bolt"
)

func TestCleanDateBuckets(t *testing.T) {
	testit(t,func (t *testing.T,db *DB) {
		maxDays := db.maxDays
		for i := 0; i < maxDays * 2; i += 1 {
			db.BagDateAdd(tnow.AddDate(0,0,i * -1),tbag,[]byte("dont"),[]byte("care"))
		}
		db.CleanDateBuckets()
		for i := 0; i < maxDays; i += 1 {

			err := db.db.View(func(tx *bolt.Tx) error {
				bname := bucketDate(tnow.AddDate(0,0,i * -1))
				assert.NotNil(t,tx.Bucket(bname))
				return nil
			})
			assert.Nil(t,err)
		}
		for i := int(maxDays+1); i < int(maxDays * 2) + 7; i += 1 {

			err := db.db.View(func(tx *bolt.Tx) error {
				bname := bucketDate(tnow.AddDate(0,0,i * -1))
				assert.Nil(t,tx.Bucket(bname))
				return nil
			})
			assert.Nil(t,err)
		}
	})
}
