/* 2019-02-14 (cc) <paul4hough@gmail.com>

Single AlertGroup Queue Manager Thread
*/
package amgr

import (
	"sync"
)

type Manager {
	c	chan
}
func NewManager() *Manager {
	return &Manager{c: make(chan bool)}
}

func (m *Manager) Notify(t time.Duration) {

	select {
	case m.c <- true:
	case <- time.After(t):
	}
}

func (am *Amgr)Manage() {

	for {
		// grab array of queue keys
		agq := a.db.AGroupQueue()

		if len(agq) < 1 {
			// wait for next alert, double check queue every 10 min
			select {
			case <- am.qmgr.c
			case <- time.After(10 * time.Minute):
			}
		}

		var wg sync.WaitGroup
		for _, agqkey := range agq {
			am.procq <- agkey
			wg.Add(1)
			go func() {
				defer wg.Done()
				am.Respond(agqkey)
			}
		}
		wg.Wait()
	}
}
