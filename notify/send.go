/* 2019-10-19 (cc) <paul4hough@gmail.com>
   send note to notify system

   lock note key
   defer unlock

   is retry?
      update retry
      return

   read db
   is new ?
      create
   else
     is close ?
       close
     else
       update
  error?
    put in retry
  else
    update db
*/
package notify

type Alert struct {
	Name	string
	Labels	pmod.LabelSet
	Annots	pmod.LabelSet
	Starts	time.Time
	Genurl	string
	Labsfp	pmod.Fingerprint
}

type Note struct {
	Labels  pmod.Labels
	Alerts  []Alert
	From	string
	Updates string
}

func (n *Notify) UpdateNote(key Key, text string) {
	// possible close b4 update - create w/ update only
	n.Send(key,Note{Updates: text})
}

func (n *Notify) Send(key Key, note Note) {

	n.Lock(key)
	defer n.UnLock(key)

	if r, ok := n.retry.Load(key); ok {
		n.retry.Store(key,note)
		return
	}

	rec := &Note

	err := n.DB(sys,grp).View(func(tx *bolt.Tx) error {
		if b := tx..Bucket(bucketName()); b == nil {
			panic( errors.NewError("note bucket not init") )
		}

		if nbuf := b.Get(key); nbuf != nil {
			if err := gob.NewDecoder(bytes.NewBuffer(nbuf)).Decode(rec); err != nil {
				panic( err )
			}
		}
		return nil
	})
	if err != nil { panic(err) }


	if nrec.nid == nil {
		nrec.nid, err = n.System(key.System).Create(note,remedCnt > 0)
	} else {
		text := note.Changes(nrec.Alerts)
		if len(note.Alerts) == 0 {
			err = n.System(key.System).Close(nrec.nid,text)
			nrec = nil
		} else {
			closed, err := n.System(key.System).Update(nrec.nid,text)
			if closed {
				// cleanup
				err := n.db.Update(func(tx *bolt.Tx) error {
					if b := tx..Bucket(bucketName()); b == nil {
						panic( errors.NewError("note bucket not init") )
					}
					return b.Delete(key)
				})
				if err != nil { panic(err) }
				nrec.nid, err = n.System(key.System).Create(note,remedCnt > 0)
			}
		}
	}
	if err != nil {
		n.retry.Store(key,note)
	} else {
		err := n.db.Update(func(tx *bolt.Tx) error {
			if b := tx..Bucket(bucketName()); b == nil {
				panic( errors.NewError("note bucket not init") )
			}
			if nrec == nil {
				return b.Delete(key)
			} else {
				var nbuf bytes.Buffer
				if err = gob.NewEncoder(&nbuf).Encode(nrec); err != nil {
					panic( err )
				} else {
					return b.Put(key,nbuf.Bytes())
				}
			}
			return nil
		})
		if err != nil { panic(err) }
	}
}
