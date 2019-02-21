/* 2019-02-13 (cc) <paul4hough@gmail.com>
   FIXME what is this for?
*/
package db

import (
	"encoding/binary"
	"errors"
	"github.com/boltdb/bolt"
)
const (
	agroupBName	= "agroup"
)

type AGroup struct {
	Json	[]byte
	Resolve	bool
}
func agroupBucket() []byte {
	return []byte(agroupBName)
}

func (db *DB) AGroupAdd(json []byte,resolve bool) {

	err := db.db.Update(func(tx *bolt.Tx) error {

		if b := tx.Bucket(agroupBucket()); b != nil {
			if key, err := b.NextSequence(); err == nil {
				keyBuf := make([]byte,binary.MaxVarintLen64)
				kn := binary.PutUvarint(keyBuf,key)

				var rbyte byte
				if resolve {
					rbyte = 1
				} else {
					rbyte = 0
				}
				db.metrics.agqueue.Inc()
				return b.Put(keyBuf[:kn],append(json,rbyte))
			} else {
				if db.debug { panic(err) }
				return err
			}
		} else {
			msg := string(agroupBucket()) + ": bucket missing"
			panic(msg)
			return errors.New(msg)
		}
	})
    if err != nil {
		panic(err)
	}
}

func (db *DB) AGroupQueue() []uint64 {

	var q []uint64

	err := db.db.View(func(tx *bolt.Tx) error {
		if b := tx.Bucket(agroupBucket()); b != nil {
			q = make([]uint64,0,b.Stats().KeyN)
			c := b.Cursor()
			for k, _ := c.First(); k != nil; k, _ = c.Next() {
				uk, _ := binary.Uvarint(k)
				q = append(q,uk)
			}
			return nil
		} else {
			msg := string(agroupBucket()) + ": bucket missing"
			panic(msg)
			return errors.New(msg)
		}
	})
    if err != nil {
		panic(err)
	}
	return q
}

func (db *DB) AGroupGet(key uint64) *AGroup  {
	ag := &AGroup{}

	err := db.db.View(func(tx *bolt.Tx) error {
		if b := tx.Bucket(agroupBucket()); b != nil {
			keyBuf := make([]byte,binary.MaxVarintLen64)
			kn := binary.PutUvarint(keyBuf,key)

			val := b.Get(keyBuf[:kn])
			if val != nil {
				ag.Resolve = uint8(val[len(val)-1]) != 0
				ag.Json = make([]byte,len(val)-1)
				copy(ag.Json,val[:len(val)-1])
			} else {
				ag = nil
			}
			return nil
		} else {
			msg := string(agroupBucket()) + ": bucket missing"
			panic(msg)
			return errors.New(msg)
		}
	})
    if err != nil {
		panic(err)
	}
	return ag
}

func (db *DB) AGroupDel(key uint64) {

	err := db.db.Update(func(tx *bolt.Tx) error {
		if b := tx.Bucket(agroupBucket()); b != nil {
			keyBuf := make([]byte,binary.MaxVarintLen64)
			kn := binary.PutUvarint(keyBuf,key)

			db.metrics.agqueue.Dec()
			return b.Delete(keyBuf[:kn])
		} else {
			msg := string(agroupBucket()) + ": bucket missing"
			panic(msg)
			return errors.New(msg)
		}
	})
	if err != nil {
		panic(err)
	}
}
