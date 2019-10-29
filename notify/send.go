/* 2019-10-19 (cc) <paul4hough@gmail.com>
   send note to notify system
*/
package notify

import (
	"fmt"
	"github.com/pahoughton/agate/notify/note"
)

func (self *Notify) Update(key Key, text string) {
	self.klock.Lock(key.KString())
	defer self.klock.Unlock(key.KString())

	note := self.dbGet(key)
	if note.Nid != nil {
		_, err := self.Sys(key.Sys).Update(note,text)
		if err != nil {
			fmt.Printf("WARN note update fail: %v for \n%s\n%s",err,note.String(),text)
		}
	} else {
		fmt.Printf("WARN note update not found %s",text)
	}
}


func (self *Notify) Send(key Key, note note.Note, remedCnt int) {

	kstr := key.KString()

	self.klock.Lock(kstr)
	defer self.klock.Unlock(kstr)

	if _, ok := self.retryMap.Load(kstr); ok {
		// replace
		self.retryMap.Store(kstr,retry{key,note,remedCnt})
		return
	}
	// functionality shared w/ retry proc
	self.notify(key,note,remedCnt)
}

func (self *Notify) notify(key Key, note note.Note, remedCnt int) {

	var err error
	rec := self.dbGet(key)
	// process
	if rec.Nid == nil {
		note.Nid, err = self.Sys(key.Sys).Create(key.Grp,note,remedCnt)
	} else {
		note.Nid = rec.Nid
		text := note.Changes(rec.Alerts)
		if len(note.Alerts) == 0 {
			err = self.Sys(key.Sys).Close(note,text)
			note.Nid = nil
		} else if len(text) > 0 {
			var closed bool
			closed, err = self.Sys(key.Sys).Update(note,text)
			if closed {
				self.dbDelete(key)
				note.Nid, err = self.Sys(key.Sys).Create(key.Grp,note,remedCnt)
			}
		}
	}
	if err != nil {
		self.retryMap.Store(key,note)
		if self.debug { fmt.Printf("DEBUG retry key: %v\n",key)	}
	} else {
		if note.Nid != nil {
			self.dbUpdate(key,note)
		} else {
			self.dbDelete(key)
		}
	}
}
