/* 2019-10-15 (cc) <paul4hough@gmail.com>
store and q note for notify sys

returns notice key
*/
package notify
/*
import (
)

type Note struct {
	key		pmod.LabelSet
	title	string
	desc	string
	labels	pmod.LabelSet
	alerts  []pmod.Fingerprint // sorted - mini btree?
	remed	uint
	resolve	bool
	nid		Nid
	updates []string // update messages
}

func (n *Notify) Queue(
	sys		string,
	grp		string,
	key		pmod.LabelSet
	title	string,
	desc	string,
	labels	pmod.LabelSet,
	alerts  []pmod.LabelSet,
	remed	uint,
	resolve	bool,
) Key {
	qkey := sys + "-" + grp

	err := n.db.Update(func(tx *bolt.Tx) error {
		if b, err := tx.CreateBucketIfNotExists(qkey); b != nil {

			keyBuf := make([]byte,binary.MaxVarintLen64)
			kl := binary.PutUvarint(keyBuf,key.Fingerprint())
			note = &Note{}
			if v := b.Get(keyBuf[:kl]); v != nil {
				dec := gob.NewDecoder(bytes.NewBuffer(v))
				if err := dec.Decode(note); err != nil {
					return err
				}
				// update must match new to old alerts and increment remed
				note.Update(alerts,remed)
			} else {
				// create
				note = &Note{
					key: key
					title: title
					desc: desc
					labels: labels
					resolve: resolve
				}
				note.Update(alerts,remed)
			}
			return b.Put(keyBuf[:kl])
		}
	})
	if err != nil {
		panic(err)
	}
	if ! n.queues[qkey] {
		go func() { n.Notify(sys,grp,n.queue[qkey]) }
	}
	n.queue[qkey] <- nkey
}
*/
