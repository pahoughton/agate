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

type AlertDB struct {
	db *bolt.DB
	maxDays int
}

const (
	BNameFmt = "2006-01-02"  // buckets named by alert date
)

func Open(dir string, mode os.FileMode, maxDays uint ) (*AlertDB, error) {
	opts := &bolt.Options{
		Timeout: 50 * time.Millisecond,
	}
	dbfn := path.Join(dir,"agate.bolt")

	bdb, err := bolt.Open(dbfn, mode, opts)
	if err != nil {
		return nil, err
	}
	adb := &AlertDB{db: bdb, maxDays: int(maxDays)}

	adb.CleanBuckets()
	// reclean alert buckets every 24 hours
	cleanBucketTicker := time.NewTicker(time.Hour * 24)
	go func() {
		for _ = range cleanBucketTicker.C {
			adb.CleanBuckets()
		}
	}()
	return adb, nil
}

func (adb *AlertDB) CleanBuckets() {

	minDate := time.Now().AddDate(0,0,adb.maxDays * -1).Format(BNameFmt)

	fmt.Println("INFO cleaning buckets before ",minDate)

	var delList []string

	err := adb.db.View(func(tx *bolt.Tx) error {

		err := tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			if strings.Compare(string(name),minDate) < 0 {
				delList = append(delList,string(name))
			}
			return nil
		})
		return err
	})
	if err != nil {
		fmt.Println("FATAL reading buckets ",err.Error())
		return
	}
	err = adb.db.Update(func(tx *bolt.Tx) error {
		for _, bname := range delList {
			if err := tx.DeleteBucket([]byte(bname)); err != nil {
				fmt.Println("ERROR delete bucket ",bname," - ",err.Error())
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println("FATAL deleting buckets ",err.Error())
		return
	}
}

func (adb *AlertDB) AddTicket(
	start	time.Time,
	fp		uint64,
	tid		string) error {

	bname := start.Format(BNameFmt)

	aKey := make([]byte, binary.MaxVarintLen64)

	binary.PutUvarint(aKey, fp)

	err := adb.db.Update(func(tx *bolt.Tx) error {

		bkt, err := tx.CreateBucketIfNotExists([]byte(bname))
		if err != nil {
			return err
		}
		return bkt.Put(aKey,[]byte(tid))
	})
	return err
}

func (adb *AlertDB) GetTicket(start time.Time,fp uint64) (string, error) {

	bname := start.Format(BNameFmt)

	aKey := make([]byte, binary.MaxVarintLen64)

	binary.PutUvarint(aKey, fp)

	var tid string

	err := adb.db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(bname))
		if bkt == nil {
			return errors.New("bucket not found " + bname)
		}
		val := bkt.Get(aKey)
		if val == nil {
			return fmt.Error("alert not found: %u", aKey)
		}
		tid = string(val)
		return nil
	})
	return tid, err
}

func (adb *AlertDB) Delete(start time.Time,fp uint64) (string, error) {

	bname := start.Format(BNameFmt)

	aKey := make([]byte, binary.MaxVarintLen64)

	binary.PutUvarint(aKey, fp)

	err := adb.db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(bname))
		if bkt == nil {
			return errors.New("bucket not found " + bname)
		}
		return bkt.Delete(aKey)
	})
	return err
}
