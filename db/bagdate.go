/* 2019-02-13 (cc) <paul4hough@gmail.com>
   FIXME what is this for?
*/
package db

import (
	"time"
	"github.com/boltdb/bolt"
	promp "github.com/prometheus/client_golang/prometheus"
)
func (db *DB) BagDateGet(t time.Time,bname,key []byte) []byte {
	var got []byte
	err := db.db.View(func(tx *bolt.Tx) error {
		if bt := tx.Bucket(bucketDate(t)); bt != nil {
			if b := bt.Bucket(bname); b != nil {
				if val := b.Get(key); val != nil {
					got = make([]byte,len(val))
					copy(got,val)
				}
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return got
}

func (db *DB) BagDateAdd(t time.Time,bname,key []byte,val []byte) {
	err := db.db.Update(func(tx *bolt.Tx) error {
		if bt, err := tx.CreateBucketIfNotExists(bucketDate(t)); bt != nil {
			if b, err := bt.CreateBucketIfNotExists(bname); b != nil {
				return b.Put(key,val)
			} else {
				panic(err)
			}
		} else {
			return err
		}
	})
	if err != nil {
		panic(err)
	}
	ml := promp.Labels{
		"date":   t.Format(TIMEFMT),
		"bucket": string(bname),
	}
	db.metrics.dbucket.With(ml).Inc()
}

func (db *DB) BagDateDel(t time.Time,bname,key []byte) {
	err := db.db.Update(func(tx *bolt.Tx) error {
		if bt := tx.Bucket(bucketDate(t)); bt != nil {
			if b := bt.Bucket(bname); b != nil {
				return b.Delete(key)
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	ml := promp.Labels{
		"date":   t.Format(TIMEFMT),
		"bucket": string(bname),
	}
	db.metrics.dbucket.With(ml).Dec()
}


func dbucketList() []string {
	// must be unique
	return []string{"alert-nid","ag-nid","nsys"}
}
func bucketAlertNid() []byte {
	return []byte(dbucketList()[0])
}
func bucketAgNid() []byte {
	return []byte(dbucketList()[1])
}
func bucketNSys() []byte {
	return []byte(dbucketList()[2])
}

func (db *DB) AlertNidGet(t time.Time,key []byte) []byte {
	return db.BagDateGet(t,bucketAlertNid(),key)
}
func (db *DB) AlertNidAdd(t time.Time,key []byte,val []byte) {
	db.BagDateAdd(t,bucketAlertNid(),key,val)
}
func (db *DB) AlertNidDel(t time.Time, key []byte) {
	db.BagDateDel(t,bucketAlertNid(),key)
}

func (db *DB) AGroupNidGet(t time.Time,key []byte) []byte {
	return db.BagDateGet(t,bucketAgNid(),key)
}
func (db *DB) AGroupNidAdd(t time.Time,key []byte,val []byte) {
	db.BagDateAdd(t,bucketAgNid(),key,val)
}
func (db *DB) AGroupNidDel(t time.Time, key []byte) {
	db.BagDateDel(t,bucketAgNid(),key)
}
