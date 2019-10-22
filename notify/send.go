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

	self.klock.Lock(key.KString())
	defer self.klock.Unlock(key.KString())

	if _, ok := self.retry.Load(key); ok {
		self.retry.Store(key,note)
		return
	}
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
		self.retry.Store(key,note)
	} else {
		if note.Nid != nil {
			self.dbUpdate(key,note)
		} else {
			self.dbDelete(key)
		}
	}
}
