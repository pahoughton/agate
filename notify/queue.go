/* 2019-10-15 (cc) <paul4hough@gmail.com>
store and q note for notify sys

returns note key
*/
package notify
import (
	"bytes"
	"encoding/gob"
	pmod "github.com/prometheus/common/model"
	promp "github.com/prometheus/client_golang/prometheus"
	"github.com/boltdb/bolt"
)

type Note struct {
	key		[]byte
	title	string
	desc	string
	labels	pmod.LabelSet
	alerts  []pmod.Fingerprint // sorted
	remed	uint
	resolve	bool
	nid		Key
	updates []string // update messages
}
func queueKey( sys, grp string) []byte {
	return []byte(sys + "-" + grp)
}

func (n *Notify) Queue(
	sys		string,
	grp		string,
	key		[]byte,
	title	string,
	desc	string,
	labels	pmod.LabelSet,
	alerts  []pmod.LabelSet,
	resolve	bool,
) Key {
	qkey := string(queueKey(sys,grp))

	nkey := &Key{Sys: sys, Grp: grp, Key: key}

	err := n.db.Update(func(tx *bolt.Tx) error {
		if b, err := tx.CreateBucketIfNotExists([]byte(qkey)); b != nil {

			note := &Note{}
			if v := b.Get(key); v != nil {
				dec := gob.NewDecoder(bytes.NewBuffer(v))
				if err := dec.Decode(note); err != nil {
					return err
				}
				// update must match new to old alerts and increment remed
				if ! note.Update(alerts) {
					nkey = nil
					return nil
				}
			} else {
				// create
				note = &Note{
					key: key,
					title: title,
					desc: desc,
					labels: labels,
					resolve: resolve,
				}
				note.Update(alerts)
			}
			var nbuf bytes.Buffer
			if err = gob.NewEncoder(&nbuf).Encode(note); err != nil {
				return err
			} else {
				return b.Put(key,nbuf.Bytes())
			}

		} else {
			return err
		}
	})
	if err != nil {
		panic(err)
	}
	if nkey != nil { // any changes

		if _, ok := n.queue[qkey]; ! ok {
			n.queue[qkey] = make(chan []byte,n.qdepth)
			go func() { n.Notify(sys,grp,n.queue[qkey]) }()
		}
		n.metrics.qlen.With(promp.Labels{"sys": sys, "grp": grp }).Inc()
		n.queue[qkey] <- nkey.Key
	}
	return *nkey
}

// update note to reflect changes in alerts
func (n *Note) Update(alerts []pmod.LabelSet) bool {
	return false
}
