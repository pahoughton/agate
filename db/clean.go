/* 2019-04-02 (cc) <paul4hough@gmail.com>
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
	TIMEFMT	= "2006-01-02"  // buckets named by alert date
)

func bucketDate(t time.Time) []byte {
	return []byte(t.Format(TIMEFMT))
}

func (db *DB) CleanDateBuckets() {

	minDate := time.Now().AddDate(0,0,db.maxDays * -1).Format(TIMEFMT)

	err := db.db.Update(func(tx *bolt.Tx) error {
		err := tx.ForEach( func(k []byte, b *bolt.Bucket) error {
			date := string(k)
			if strings.Compare(date,minDate) < 0 {
				fmt.Println("INFO remove bucket ",date)
				for _, n := range dbucketList() {
					ml := promp.Labels{"date": date,"bucket": n}
					db.metrics.dbucket.Delete(ml)
				}
				tx.DeleteBucket(k)
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
