/* 2018-12-27 (cc) <paul4hough@gmail.com>
   notification system interface
*/
package notify

import (
    "encoding/base64"
	"fmt"
	"sync"

	"github.com/pahoughton/agate/config"
	"github.com/pahoughton/agate/notify/note"
//	"github.com/pahoughton/agate/notify/mock"
//	"github.com/pahoughton/agate/notify/gitlab"
	"github.com/pahoughton/agate/notify/hpsm"

	proma "github.com/prometheus/client_golang/prometheus/promauto"
	promp "github.com/prometheus/client_golang/prometheus"

	"github.com/xiaonanln/keylock"
)

const (
	SysMock	= "mock"
	SysGitlab	= "gitlab"
	SysHpsm	= "hpsm"
)


type System interface {
	Create(goup string, note note.Note, remcnt int) ([]byte, error)
	Update(note note.Note, text string) (bool,error)
	Close(note note.Note, text string) error
	Group() string
	Name() string
}

type metrics struct {
	notes		*promp.CounterVec
	errors		promp.Counter
}

type Notify struct {
	DefSys			string
	CloseResolved	bool
	dataDir			string
	db				*DB
	sys				map[string]System
	klock			keylock.KeyLock
	retry			sync.Map
	metrics			metrics
	debug			bool
}

type Key struct {
	Sys		string
	Grp		string
	Key		[]byte
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
		dataDir:		dataDir,
		db:				newDB(),
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
	self.sys = make(map[string]System,16)
	self.sys[SysHpsm] = hpsm.New(SysHpsm, cfg.Sys.Hpsm,dbg)

	if _, ok  := self.sys[self.DefSys]; ! ok {
		self.unregister()
		panic(fmt.Sprintf("invalid default ticket sys: %s",cfg.Default))
	}

	return self
}

func (n *Notify) Del() {
	n.unregister()
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

func (n *Notify) Sys(sys string) System {
	if s, ok := n.sys[sys]; ok {
		return s
	} else {
		return nil
	}
}

func (n *Notify) Group(nsys string) string {
	if n.Sys(nsys) != nil {
		return n.Sys(nsys).Group()
	} else {
		return "invalid"
	}
}
