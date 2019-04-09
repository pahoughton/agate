/* 2019-02-14 (cc) <paul4hough@gmail.com>

Single AlertGroup Queue Manager Thread
*/
package amgr

import (
	"os"
	"time"
	"github.com/pahoughton/agate/notify"
)
type NSysChan struct {
	c	chan bool
}
type Manager struct {
	nsys	[]NSysChan
	quit	chan bool
}
func NewManager() *Manager {
	m := &Manager{
		nsys: make([]NSysChan,notify.NSysUnknown),
		quit: make(chan bool),
	}
	for i := notify.NSysMock; i < notify.NSysUnknown; i += 1 {
		m.nsys[i].c = make(chan bool)
	}
	return m
}

func (m *Manager) Notify(nsys uint) {
	select {
	case m.nsys[nsys].c <- true:
	case <- time.After(1):
	}
}

func (m *Manager) Quit() {
	select {
	case m.quit <- true:
	case <- time.After(1):
	}
}

func (am *Amgr) manage(nsys uint) {

	for {
		// grab array of queue keys
		agq := am.db.AGroupQueueList(nsys)
		if len(agq) < 1 {
			// wait for next alert, double check queue every 10 min
			select {
			case <- am.qmgr.nsys[nsys].c:
			case <- am.qmgr.quit:
				os.Exit(0)
			case <- time.After(am.retry):
			}
		} else {
			for _, id := range agq {
				if am.Notify(notify.NSys(nsys),id) == false {
					time.Sleep(am.retry)
					break;
				}
			}
		}
	}
}

func (am *Amgr) Manage() {

	go func() { am.manage(uint(notify.NSysMock)) }()
	go func() {	am.manage(uint(notify.NSysGitlab)) }()
	go func() {	am.manage(uint(notify.NSysHpsm)) }()
}
