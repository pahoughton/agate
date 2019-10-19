/* 2019-09-03 (cc) <paul4hough@gmail.com>
*/
package notify

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"fmt"
	"strconv"
	"time"
	"github.com/boltdb/bolt"
	promp "github.com/prometheus/client_golang/prometheus"
)

type Notice struct {
	key			[]byte
	resolved	bool
	title		string
	note		string
}

// hackish - name for metric - order from notify.NSys notify/new.go
func nsysname(s uint) string {
	name := []string{"mock","gitlab","hpsm"}
	if int(s) < len(name) {
		return name[s] + "-" +  strconv.Itoa(int(s))
	} else {
		return strconv.Itoa(int(s))
	}
}

func qbucket(nsys uint, ngrp string) []byte {
	return []byte(fmt.Sprintf("nq-%s-%s",nsysname(nsys),ngrp))
}
func kbucket(nsys uint, ngrp string) []byte {
	return []byte(fmt.Sprintf("nk-%s-%s",nsysname(nsys),ngrp))
}

func (notify *N) Queue(
	nsys		uint,
	ngrp		string,
	nkey		[]byte,
	resolved	bool,
	title		string,
	note		string) {

	err := n.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(qbucket(nsys,ngrp))
		if err != nil {
			return err
		}
		var note Notice{nkey, resolved, title, note}

		var datBuf bytes.Buffer
		if err = gob.NewEncoder(&datBuf).Encode(note); err != nil {
			return err
		}
		if seq, err := b.NextSequence(); err != nil {
			return err
		}
		keyBuf := make([]byte,binary.MaxVarintLen64)
		kn := binary.PutUvarint(keyBuf,seq)

		return b.Put(keyBuf[:kn],datbuf)

	})
	if err != nil {
		panic(err)
	}
}

func (notify *N) NidAdd(nsys uint, ngrp string, nkey, nid []byte) {
	err := n.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(kbucket(nsys,ngrp))
		if err != nil {
			return err
		}

		return b.Put(nkey,nid)
	})
	if err != nil {
		panic(err)
	}
}

func (notify *N) NidGet(nsys uint, ngrp string, nkey []byte) []byte {
	var nid []byte
	err := n.db.View(func(tx *bolt.Tx) error {
		b, err := tx.Bucket(kbucket(nsys,ngrp)); err != nil {
			return nil
		}

		if val := b.Get(nBuf); val != nil {
			nid = make([]byte,len(val))
			copy(nid,val)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return nid
}







func (db *DB) AGroupNSysDel(t time.Time,key []byte) {
	db.BagDateDel(t,bucketNSys(),key)
}

func (db *DB) AGroupNSysGet(t time.Time,key []byte) *NSys {

	var got *NSys

	err := db.db.View(func(tx *bolt.Tx) error {
		if bt := tx.Bucket(bucketDate(t)); bt != nil {
			if b := bt.Bucket(bucketNSys()); b != nil {
				if val := b.Get(key); val != nil {
					dec := gob.NewDecoder(bytes.NewBuffer(val))
					got = &NSys{}
					return dec.Decode(got)
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

func (db *DB) AGroupQueueNSysAdd(t time.Time,nsys NSys,agkey,agval []byte) {

	err := db.db.Update(func(tx *bolt.Tx) error {
		bt, err := tx.CreateBucketIfNotExists(bucketDate(t))
		if bt != nil {
			if b, err := bt.CreateBucketIfNotExists(bucketNSys()); b != nil {
				var val bytes.Buffer
				enc := gob.NewEncoder(&val)
				if err = enc.Encode(nsys); err == nil {
					if err = b.Put(agkey,val.Bytes()); err != nil {
						return err
					} else {
						ml := promp.Labels{
							"date":   t.Format(TIMEFMT),
							"bucket": string(bucketNSys()),
						}
						db.metrics.dbucket.With(ml).Inc()
					}
				} else {
					return err
				}
			} else {
				return err
			}
		} else {
			return err
		}

		if b, err := tx.CreateBucketIfNotExists(bucketAg(nsys.Sys)); b != nil {
			if seq, err := b.NextSequence(); err == nil {

				ml := promp.Labels{"sys": nsysname(nsys.Sys)}
				db.metrics.agqueue.With(ml).Inc()

				keyBuf := make([]byte,binary.MaxVarintLen64)
				kn := binary.PutUvarint(keyBuf,seq)

				return b.Put(keyBuf[:kn],agval)
			} else {
				return err
			}
		} else {
			return err
		}
	})
	if err != nil {
		panic(err)
	}
}

func (db *DB) AGroupQueueAdd(nsys uint,agval []byte) {

	err := db.db.Update(func(tx *bolt.Tx) error {
		if b, err := tx.CreateBucketIfNotExists(bucketAg(nsys)); b != nil {
			if seq, err := b.NextSequence(); err == nil {

				ml := promp.Labels{"sys": nsysname(nsys)}
				db.metrics.agqueue.With(ml).Inc()

				keyBuf := make([]byte,binary.MaxVarintLen64)
				kn := binary.PutUvarint(keyBuf,seq)

				return b.Put(keyBuf[:kn],agval)
			} else {
				return err
			}
		} else {
			return err
		}
	})
	if err != nil {
		panic(err)
	}
}

func (db *DB) AGroupQueueList(nsys uint) []uint64 {

	q := make([]uint64,0,16)
	err := db.db.View(func(tx *bolt.Tx) error {
		if b := tx.Bucket(bucketAg(nsys)); b != nil {
			q = make([]uint64,0,b.Stats().KeyN)
			c := b.Cursor()
			for k, _ := c.First(); k != nil; k, _ = c.Next() {
				uk, _ := binary.Uvarint(k)
				q = append(q,uk)
			}
		}
		return nil
	})
    if err != nil {
		panic(err)
	}
	return q
}

func (db *DB) AGroupQueueGet(nsys uint, key uint64) []byte  {

	var ag []byte
	err := db.db.View(func(tx *bolt.Tx) error {
		if b := tx.Bucket(bucketAg(nsys)); b != nil {
			keyBuf := make([]byte,binary.MaxVarintLen64)
			kn := binary.PutUvarint(keyBuf,key)

			if val := b.Get(keyBuf[:kn]); val != nil {
				ag = make([]byte,len(val))
				copy(ag,val)
			}
		} else {
			msg := string(bucketAg(nsys)) + ": bucket missing"
			panic(msg)
			return errors.New(msg)
		}
		return nil
	})
    if err != nil {
		panic(err)
	}
	return ag
}

func (db *DB) AGroupQueueDel(nsys uint, key uint64) {

	err := db.db.Update(func(tx *bolt.Tx) error {
		if b := tx.Bucket(bucketAg(nsys)); b != nil {
			keyBuf := make([]byte,binary.MaxVarintLen64)
			kn := binary.PutUvarint(keyBuf,key)

			ml := promp.Labels{"sys": nsysname(nsys)}
			db.metrics.agqueue.With(ml).Dec()

			if err := b.Delete(keyBuf[:kn]); err != nil {
				if b.Stats().KeyN == 0 {
					if err = b.SetSequence(0); err != nil {
						return err
					}
				}
				return nil
			} else {
				return err
			}
		} else {
			msg := string(bucketAg(nsys)) + ": bucket missing"
			panic(msg)
			return errors.New(msg)
		}
	})
	if err != nil {
		panic(err)
	}
}
