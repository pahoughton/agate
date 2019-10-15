/* 2019-10-15 (cc) <paul4hough@gmail.com>
   FIXME what is this for?
*/
package remed

import (
	"testing"
)

const (
	taskName = "alertname"
	bucketName = "remed"
)

// todo - what if never receive resolve?
// remediated & unresolved metric
func (r *Remed) Queue(n pmod.LabelSet, nkey []byte, resolved bool) {

	if ! n[taskName] || ! r.Has(n[taskName]) {
		return
	}

	err := db.db.Update(func(tx *bolt.Tx) error {

		if b := tx.Bucket(bucketName); b != nil {

			keyBuf := make([]byte,binary.MaxVarintLen64)
			kl := binary.PutUvarint(keyBuf,n.Fingerprint())
			if v := b.Get(keyBuf[:kl]); v != nil {
				// remediation has been fired
				if resolved {
					r.metrics.unres.Dec()
					return b.Del(keyBuf[:kl])
				}
			} else if ! resolved {
				// new unresolved
				if err := b.Put(keyBuf[:kl]); err == nil {
					r.metrics.unres.Inc()
					return r.Do(n,nkey)
				}
			}
		} else {
			panic( "remed queue not initialized" )
		}
	})
	if err != nil {
		panic(err)
	}
}
