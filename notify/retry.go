/* 2019-10-22 (cc) <paul4hough@gmail.com>
   FIXME what is this for?
*/
package notify

import "time"

// todo metrics & validation
func (self *Notify) retry() {

	for {
		self.retryOnce()
		time.Sleep(self.RetryDelay)
	}
}
// broke out for testing
func (self *Notify) retryOnce() {

	tried := make(map[string]bool)
	for {
		// grab one
		var kstr string
		var r retry

		self.retryMap.Range(func(k,v interface{}) bool {
			kstr = k.(string)
			if _, ok := tried[kstr]; ! ok {
				r = v.(retry)
				return false
			}
			return true
		})
		// found one?
		if len(r.key.Sys) > 0 {
				tried[kstr] = true
			self.klock.Lock(kstr)
			self.retryMap.Delete(kstr)
			self.notify(r.key,r.note,r.rcnt)
			self.klock.Unlock(kstr)
		} else {
			break
		}
	}
}
