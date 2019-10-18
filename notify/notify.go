/* 2019-02-19 (cc) <paul4hough@gmail.com>
*/
package notify

import (
	"errors"
	"bytes"
	"encoding/gob"
	/*
	pmod "github.com/prometheus/common/model"
	promp "github.com/prometheus/client_golang/prometheus"
*/
	"github.com/boltdb/bolt"
)

func (n *Notify) Notify(sys, grp string, q chan []byte) {
	for {
		key := <- q
		err := n.db.Update(func(tx *bolt.Tx) error {
			if b := tx.Bucket(queueKey(sys,grp)); b != nil {

				note := &Note{}
				if nbuf := b.Get(key); nbuf != nil {
					dec := gob.NewDecoder(bytes.NewBuffer(nbuf))
					if err := dec.Decode(note); err != nil {
						return err
					} else {
						if err = n.notify(sys,grp,*note); err == nil {
							var buf bytes.Buffer
							err = gob.NewEncoder(&buf).Encode(note)
							if err != nil {
								return err
							} else {
								return b.Put(key,buf.Bytes())
							}
						} else {
							return err
						}
					}
				} else {
					panic("FIXME metric for key not found - not error")
				}
			} else {
				return errors.New("notify init error")
			}
		})
		if err != nil {
			panic(err)
		}
	}
}

func (n *Notify) notify(sys, grp string, note Note) error {
	panic("FIXME STUB")
	return nil
}
/*
import (
	"fmt"
	promp "github.com/prometheus/client_golang/prometheus"
)



func (n *Notify) Create(
	nsys	string,
	grp		string,
	title	string,
	desc	string,
	remed	bool,
	resolve	bool) (Key, error) {

	if n.System(nsys) != nil {
		var (
			aclose string
			aremed string
		)
		if resolve {
			aclose = "closes on resolve"
		} else {
			aclose = "manual"
		}
		if remed {
			aremed = "true"
		} else {
			aremed = "false"
		}
		ndesc := fmt.Sprintf(
			"\nauto-close: %s  remediation: %s\n%s",
			aclose,
			aremed,
			desc)

		nid, err := n.System(nsys).Create(grp,title,ndesc)
		if err == nil {
			n.metrics.notes.With(promp.Labels{
				"sys": nsys.String(),
				"grp": grp,
			}).Inc()
			return nid, err
		} else {
			n.metrics.errors.Inc()
			return nid, err
		}
	} else {
		panic(fmt.Sprintf("invalid nsys: %d\n",nsys))
		return nil, nil
	}
}

func (n *Notify) Update(nid Key, msg string) bool {
	if n.System(string(nid.Sys())) != nil {
		err := n.System(string(nid.Sys())).Update(nid,msg)
		if err == nil {
			return true
		} else {
			n.metrics.errors.Inc()
		}
	}
	return false
}

func (n *Notify) Close(nid Key, msg string) bool {
	if n.System(string(nid.Sys())) != nil {
		err := n.System(string(nid.Sys())).Close(nid,msg)
		if err == nil {
			return true
		} else {
			n.metrics.errors.Inc()
		}
	}
	return false
}
*/
