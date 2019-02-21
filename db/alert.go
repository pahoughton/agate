/* 2019-02-13 (cc) <paul4hough@gmail.com>
   FIXME what is this for?
*/
package db

import (
	"fmt"
	"strings"
	"time"
	"github.com/boltdb/bolt"
	promp "github.com/prometheus/client_golang/prometheus"
)

const (
	BNameFmt	= "2006-01-02"  // buckets named by alert date
	alertsBName = "alerts"
)

func alertsBucket() []byte {
	return []byte(alertsBName)
}

func startBucket(start time.Time) []byte {
	return []byte(start.Format(BNameFmt))
}

func (db *DB) AlertCleanBuckets() {

	minDate := time.Now().AddDate(0,0,db.maxDays * -1).Format(BNameFmt)

	fmt.Println("INFO cleaning buckets before ",minDate)

	err := db.db.Update(func(tx *bolt.Tx) error {
		ab := tx.Bucket(alertsBucket())
		if ab == nil {
			return nil
		}
		err := ab.ForEach( func(k, v []byte) error {
			if v == nil {
				date := string(k)
				if strings.Compare(date,minDate) < 0 {
					fmt.Println("INFO remove bucket ",date)
					ab.DeleteBucket(k)
					mlabels := promp.Labels{"date": date}
					db.metrics.tickets.Delete(mlabels)
				}
			}
			return nil
		})
		return err
	})
	if err != nil {
		db.metrics.errors.Inc()
		fmt.Println("ERROR clean buckets ",err.Error())
		if db.debug { panic(err); }
	}
}


func (db *DB) AlertAdd(start time.Time,key,val []byte) {

	err := db.db.Update(func(tx *bolt.Tx) error {
		bname := startBucket(start)
		if ab := tx.Bucket(alertsBucket()); ab != nil {
			if b, err := ab.CreateBucketIfNotExists(bname); b != nil {
				db.metrics.tickets.With(
					promp.Labels{"date":string(bname)},
				).Inc()
				return b.Put(key,val)
			} else {
				panic(err)
			}
		} else {
			panic(string(alertsBucket()))
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func (db *DB) AlertGet(start time.Time,key []byte) []byte {

	var got []byte
	err := db.db.View(func(tx *bolt.Tx) error {
		bname := startBucket(start)
		if ab := tx.Bucket(alertsBucket()); ab != nil {
			if b := ab.Bucket(bname); b != nil {
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

func (db *DB) AlertDel(start time.Time,key []byte) {

	err := db.db.Update(func(tx *bolt.Tx) error {
		bname := startBucket(start)
		if ab := tx.Bucket(alertsBucket()); ab != nil {
			if b := ab.Bucket(bname); b != nil {
				db.metrics.tickets.With(
					promp.Labels{"date":string(bname)},
				).Dec()
				return b.Delete(key)
			} // else {	panic(string(bname) + " bucket missing") }
			// although surprising, no real harm
		} else {
			panic(string(alertsBucket()))
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}
