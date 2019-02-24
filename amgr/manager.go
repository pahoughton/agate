/* 2019-02-14 (cc) <paul4hough@gmail.com>

Single AlertGroup Queue Manager Thread
*/
package amgr

import (
	"os"
	"time"
)
type Manager struct {
	c	chan bool
	q	chan bool
}
func NewManager() *Manager {
	return &Manager{c: make(chan bool), q: make(chan bool)}
}

func (m *Manager) Notify(t time.Duration) {

	select {
	case m.c <- true:
	case <- time.After(t):
	}
}

func (m *Manager) Quit() {
	select {
	case m.q <- true:
	case <- time.After(1):
	}
}

func (am *Amgr) Manage() {

	for {
		// grab array of queue keys
		agq := am.db.AGroupQueue()
		if len(agq) < 1 {
			// wait for next alert, double check queue every 10 min
			select {
			case <- am.qmgr.c:
			case <- am.qmgr.q:
				os.Exit(0)
			case <- time.After(am.retry):
			}
		} else {
			for _, id := range agq {
				if am.Respond(id) == false {
					time.Sleep(am.retry)
					break;
				}
			}
		}
	}
}
