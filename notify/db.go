/* 2019-10-21 (cc) <paul4hough@gmail.com>
   agate notify persistent data
*/
package notify

import (
	"bytes"
	"errors"
	"encoding/gob"
	"path"
	"sync"
	"time"
	"github.com/pahoughton/agate/notify/note"
	"github.com/boltdb/bolt"
)

type DB struct {
	mutex	sync.Mutex
	dbmap	map[string]*bolt.DB
}

func newDB() (*DB) {
	db := &DB{}
	db.dbmap = make(map[string]*bolt.DB,16)
	return db
}

func (self *Notify) DB(sys, grp string) *bolt.DB {

	grpFn := sys + "-" + grp + "-queue.bolt"

	self.db.mutex.Lock()
	defer self.db.mutex.Unlock()


	if db, ok := self.db.dbmap[grpFn]; ! ok {
		fn := path.Join(self.dataDir,grpFn)
		opts := &bolt.Options{Timeout: 50 * time.Millisecond}
		db, err := bolt.Open(fn,0664,opts)
		if err != nil {
			panic( err )
		}
		self.db.dbmap[grpFn] = db
		err = db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists(bucketName())
			return err
		})
		if err != nil {
			panic(err)
		}
		return db
	} else {
		return db
	}
}

func (self *Notify) dbGet(key Key) (note.Note) {
	// get
	rec := note.Note{}

	err := self.DB(key.Sys,key.Grp).View(func(tx *bolt.Tx) error {
		if b := tx.Bucket(bucketName()); b == nil {
			panic( errors.New("note bucket not init") )
		} else {

			if nbuf := b.Get(key.Key); nbuf != nil {
				err := gob.NewDecoder(bytes.NewBuffer(nbuf)).Decode(rec)
				if err != nil {
					panic( err )
				}
			}
		}
		return nil
	})
	if err != nil { panic(err) }
	return rec
}

func (self *Notify) dbUpdate(key Key, note note.Note) {

	err := self.DB(key.Sys,key.Grp).Update(func(tx *bolt.Tx) error {
		if b := tx.Bucket(bucketName()); b == nil {
			panic( errors.New("note bucket not init") )
		} else {
			var nbuf bytes.Buffer
			if err := gob.NewEncoder(&nbuf).Encode(note); err != nil {
				panic( err )
			} else {
				return b.Put(key.Key,nbuf.Bytes())
			}
		}
		return nil
	})
	if err != nil { panic(err) }
}

func (self *Notify) dbDelete(key Key) {
	err := self.DB(key.Sys,key.Grp).Update(func(tx *bolt.Tx) error {
		if b := tx.Bucket(bucketName()); b == nil {
			panic( errors.New("note bucket not init") )
		} else {
			return b.Delete(key.Key)
		}
	})
	if err != nil { panic(err) }
}
