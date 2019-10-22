/* 2019-10-15 (cc) <paul4hough@gmail.com>
   FIXME what is this for?
*/
package remed

import (
	"errors"
	"encoding/binary"
	pmod "github.com/prometheus/common/model"
	"github.com/boltdb/bolt"
	"github.com/pahoughton/agate/notify"
)


// todo - what if never receive resolve?
// remediated & unresolved metric
func (r *Remed) Queue(labels pmod.LabelSet, key notify.Key, resolved bool) {

	if task, ok := labels[taskName]; ! ok  || ! r.TaskHasRemed(string(task)) {
		return
	}

	bkey := make([]byte,binary.MaxVarintLen64)
	kl := binary.PutUvarint(bkey,uint64(labels.Fingerprint()))

	err := r.db.Update(func(tx *bolt.Tx) error {

		if b := tx.Bucket([]byte(bucketName)); b != nil {

			if v := b.Get(bkey[:kl]); v != nil {
				// remediation has been fired
				if resolved {
					r.metrics.unres.Dec()
					return b.Delete(bkey[:kl])
				} else {
					return nil
				}
			} else if ! resolved {
				// new unresolved
				if err := b.Put(bkey[:kl],[]byte("rem")); err == nil {
					r.metrics.unres.Inc()
					r.Remed(string(labels[pmod.LabelName(taskName)]),labels,key)
					return nil
				} else {
					return err
				}
			} else {
				return nil
			}
		} else {
			panic( "remed queue not initialized" )
			return errors.New("remed q")
		}

	})
	if err != nil {
		panic(err)
	}
}
