/* 2018-12-27 (cc) <paul4hough@gmail.com>
   notification system interface
*/
package notify

import (
	"fmt"

	"github.com/pahoughton/agate/config"
	"github.com/pahoughton/agate/notify/mock"
	"github.com/pahoughton/agate/notify/gitlab"
	"github.com/pahoughton/agate/notify/hpsm"

	proma "github.com/prometheus/client_golang/prometheus/promauto"
	promp "github.com/prometheus/client_golang/prometheus"

	"github.com/boltdb/bolt"
)


type Key struct {
	Sys		string
	Grp		string
	Key		[]byte
}

func (n *Notify) UpdateNote(k Key, text string) bool {
	print("STUB")
	return false
}

type System interface {
	Create(goup, title, desc string) ([]byte, error)
	Update(nid []byte, desc string) error
	Close(nid []byte, desc string) error
	Group() string
}

const (
	NSysMock	= "mock"
	NSysGitlab	= "gitlab"
	NSysHpsm	= "hpsm"
)

type metrics struct {
	notes		*promp.CounterVec
	errors		promp.Counter
	qlen		*promp.GaugeVec
}

type Notify struct {
	DefSys			string
	CloseResolved	bool
	metrics			metrics
	systems			map[string]System
	queue			map[string]chan []byte
	qdepth			uint
	db				*bolt.DB
	debug			bool
}

func New(cfg config.Notify, db *bolt.DB, dbg bool) *Notify {

	n := &Notify{
		debug:			dbg,
		DefSys:			cfg.Default,
		CloseResolved:	cfg.Resolved,
		db:				db,
		qdepth:			cfg.NQDepth,
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
			qlen:	proma.NewGaugeVec(
				promp.GaugeOpts{
					Namespace:	"agate",
					Subsystem:	"notify",
					Name:		"qlen",
					Help:		"notify queue len",
				},[]string{"sys","grp"}),
		},
	}
	n.systems = make(map[string]System)
	n.systems[NSysMock]	= mock.New(cfg.Sys.Mock,NSysMock,dbg)
	n.systems[NSysGitlab] = gitlab.New(cfg.Sys.Gitlab,NSysGitlab,dbg)
	n.systems[NSysHpsm]	= hpsm.New(cfg.Sys.Hpsm,NSysHpsm,dbg)

	if _, ok  := n.systems[n.DefSys]; ! ok {
		n.unregister()
		panic(fmt.Sprintf("invalid default ticket sys: %s",cfg.Default))
	}

	return n
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

func (n *Notify) System(nsys string) System {
	if s, ok := n.systems[nsys]; ok {
		return s
	} else {
		return nil
	}
}

func (n *Notify) Group(nsys string) string {
	if n.System(nsys) != nil {
		return n.System(nsys).Group()
	} else {
		return "invalid"
	}
}
