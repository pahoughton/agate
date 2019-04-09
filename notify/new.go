/* 2018-12-27 (cc) <paul4hough@gmail.com>
   notification system interface
*/
package notify

import (
	"fmt"

	"github.com/pahoughton/agate/config"
	"github.com/pahoughton/agate/notify/nid"
	"github.com/pahoughton/agate/notify/mock"
	"github.com/pahoughton/agate/notify/gitlab"
	"github.com/pahoughton/agate/notify/hpsm"

	proma "github.com/prometheus/client_golang/prometheus/promauto"
	promp "github.com/prometheus/client_golang/prometheus"
)

type System interface {
	Create(goup, title, desc string) (nid.Nid, error)
	Update(nid nid.Nid, desc string) error
	Close(nid nid.Nid, desc string) error
	Group() string
}

type NSys int
const (
	NSysMock	NSys = iota
	NSysGitlab
	NSysHpsm
	NSysUnknown
)

type NSysMap map[string]NSys
var (
	nsysnames = []string{"mock","gitlab","hpsm","unknown"}
	nsysmap = NSysMap{
		"mock":		NSysMock,
		"gitlab":	NSysGitlab,
		"hpsm":		NSysHpsm,
		"unknown":	NSysUnknown,
	}
)

type metrics struct {
	notes		*promp.CounterVec
	errors		promp.Counter
}

type Notify struct {
	DefSys			NSys
	CloseResolved	bool
	metrics			metrics
	systems			[]System
	debug			bool
}

func NewNSys(s string) NSys {
	if v, ok := nsysmap[s]; ok {
		return v
	} else {
		return NSysUnknown
	}
}

func (t NSys) Int() int {
	return int(t)
}

func (t NSys) String() string {
	if NSysMock <= t && t <= NSysUnknown {
		return nsysnames[t]
	} else {
		return "invalid"
	}
}

func New(cfg config.Notify, dbg bool) *Notify {

	n := &Notify{
		debug:			dbg,
		DefSys:			NewNSys(cfg.Default),
		CloseResolved:	cfg.Resolved,
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
	n.systems = make([]System,len(nsysmap))
	n.systems[NSysMock]	= mock.New(cfg.Sys.Mock,NSysMock.Int(),dbg)
	n.systems[NSysGitlab] = gitlab.New(cfg.Sys.Gitlab,NSysGitlab.Int(),dbg)
	n.systems[NSysHpsm]	= hpsm.New(cfg.Sys.Hpsm,NSysHpsm.Int(),dbg)

	if NSysMock > n.DefSys || n.DefSys >= NSysUnknown {
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

func (n *Notify) System(nsys NSys) System {
	if NSysMock <= nsys && nsys < NSysUnknown {
		return n.systems[nsys]
	} else {
		return nil
	}
}
