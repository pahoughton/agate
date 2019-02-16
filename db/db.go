/* 2018-12-31 (cc) <paul4hough@gmail.com>
   agate alert db stores ticket id's for alerts so that tickets for resolved
   alerts can be updated.

   Each date will have it's own buckets to provide for deleting
   unresolved alerts that are older than 'MaxAge' days

*/
package db

import (
	"errors"
	"encoding/binary"
	"fmt"
	"os"
	"path"
	"strings"
	"time"


	"github.com/boltdb/bolt"
)

type AgateDB struct {
	db *bolt.DB
	maxDays int
}

func Open(dir string, mode os.FileMode, maxDays uint ) (*AgateDB, error) {
	opts := &bolt.Options{
		Timeout: 50 * time.Millisecond,
	}
	dbfn := path.Join(dir,"agate.bolt")

	bdb, err := bolt.Open(dbfn, mode, opts)
	if err != nil {
		return nil, fmt.Errorf("open %s %v - %v",dbfn,mode,err)
	}
	adb := &AgateDB{db: bdb, maxDays: int(maxDays)}

	adb.TicketCleanBuckets()
	// reclean alert buckets every 24 hours
	cleanBucketTicker := time.NewTicker(time.Hour * 24)
	go func() {
		for _ = range cleanBucketTicker.C {
			adb.TicketCleanBuckets()
		}
	}()
	return adb, nil
}
