/* 2018-12-27 (cc) <paul4hough@gmail.com>
   notification system interface
*/
package notify

import (
    "encoding/base64"
	"fmt"
	"sync"
	"time"

	"github.com/pahoughton/agate/config"
	"github.com/pahoughton/agate/notify/note"
	"github.com/pahoughton/agate/notify/mock"
	"github.com/pahoughton/agate/notify/gitlab"
	"github.com/pahoughton/agate/notify/hpsm"

	proma "github.com/prometheus/client_golang/prometheus/promauto"
	promp "github.com/prometheus/client_golang/prometheus"

	"github.com/xiaonanln/keylock"
)

type Key struct {
	Sys		string
	Grp		string
	Key		[]byte
}

// if change update ValidSys() - todo sys hard code
const (
	SysMock		= "mock"
	SysGitlab	= "gitlab"
	SysHpsm		= "hpsm"
)

type System interface {
	Create(group string, note note.Note, remcnt int) ([]byte, error)
	Update(note note.Note, text string) (bool,error)
	Close(note note.Note, text string) error
	Group() string
	Name() string
}

type metrics struct {
	notes		*promp.CounterVec
	errors		promp.Counter
}

type retry struct {
	key Key
	note note.Note
	rcnt int
}
type Notify struct {
	DefSys			string
	CloseResolved	bool
	RetryDelay		time.Duration
	dataDir			string
	db				*DB
	sys				map[string]System
	klock			*keylock.KeyLock
	retryMap		sync.Map
	metrics			metrics
	debug			bool
}

func bucketName() []byte { return []byte("notes"); }

func (self *Key) KString() string {
	return base64.StdEncoding.EncodeToString(self.Key)
}

func New(cfg config.Notify, dataDir string, dbg bool) *Notify {

	self := &Notify{
		debug:			dbg,
		DefSys:			cfg.Default,
		CloseResolved:	cfg.Resolved,
		RetryDelay:		cfg.Retry,
		dataDir:		dataDir,
		db:				newDB(),
		sys:			make(map[string]System,16),
		klock:			keylock.NewKeyLock(),
		metrics:		metrics{
			notes:	proma.NewCounterVec(
				promp.CounterOpts{
					Namespace:	"agate",
					Subsystem:	"notes",
					Name:		"generated",
					Help:		"number of ticekts created",
				}, []string{
					"sys",
					"grp",
				}),
			errors: proma.NewCounter(
				promp.CounterOpts{
					Namespace:	"agate",
					Subsystem:	"notify",
					Name:		"errors",
					Help:		"number of ticket errors",
				}),
		},
	}
	self.sys[SysMock] = mock.New(SysMock, cfg.Sys.Mock,dbg)
	self.sys[SysGitlab] = gitlab.New(SysGitlab, cfg.Sys.Gitlab,dbg)
	self.sys[SysHpsm] = hpsm.New(SysHpsm, cfg.Sys.Hpsm,dbg)

	if _, ok  := self.sys[self.DefSys]; ! ok {
		self.unregister()
		panic(fmt.Sprintf("invalid default ticket sys: %s",cfg.Default))
	}

	return self
}

func (self *Notify) Del() {
	self.unregister()
	for _, v := range self.db.dbmap {
		v.Close()
	}
}
func (n *Notify) unregister() {
	if n != nil &&  n.metrics.errors != nil {
		promp.Unregister(n.metrics.errors);
		n.metrics.errors = nil
	}
	if n.metrics.notes != nil  {
		promp.Unregister(n.metrics.notes);
		n.metrics.notes = nil
	}
}

func (n *Notify) Errorf(format string, args ...interface{}) error {
	n.metrics.errors.Inc()
	fmt.Println("ERROR: ",fmt.Sprintf(format,args...))
	return fmt.Errorf(format,args...)
}

func (self *Notify) ValidSys(sys string) bool {
	if len(self.sys) > 0 {
		_, ok := self.sys[sys]
		return ok
	} else {
		// todo sys hard code
		for _, s := range []string{SysMock,SysGitlab,SysHpsm} {
			if sys == s { return true }
		}
		return false
	}
}

func (self *Notify) Sys(sys string) System {
	if self.ValidSys(sys) {
		return self.sys[sys]
	} else {
		return nil
	}
}

func (n *Notify) Group(nsys string) string {
	if n.Sys(nsys) != nil {
		return n.Sys(nsys).Group()
	} else {
		return ""
	}
}
