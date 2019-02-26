/* 2018-12-27 (cc) <paul4hough@gmail.com>
   ticket management interface
*/
package ticket

import (
	"fmt"

	"github.com/pahoughton/agate/config"
	"github.com/pahoughton/agate/ticket/gitlab"
	"github.com/pahoughton/agate/ticket/hpsm"
	"github.com/pahoughton/agate/ticket/mock"
	"github.com/pahoughton/agate/ticket/tid"

	proma "github.com/prometheus/client_golang/prometheus/promauto"
	promp "github.com/prometheus/client_golang/prometheus"
)


type TSys int
const (
	TSysMock	TSys = iota
	TSysGitlab
	TSysHpsm
	TSysUnknown
)

type TSysMap map[string]TSys
var (
	tsysmap = TSysMap{
		"mock":		TSysMock,
		"gitlab":	TSysGitlab,
		"hpsm":		TSysHpsm,
		"unknown":	TSysUnknown,
	}
)

func NewTSys(s string) TSys {
	if v, ok := tsysmap[s]; ok {
		return v
	} else {
		return TSysUnknown
	}
}
func (t TSys) Int() int {
	return int(t)
}
func (t TSys) String() string {

	names := []string{"unk","gitlab","hpsm","mock"}
	if TSysMock <= t && t <= TSysUnknown {
		return names[t]
	} else {
		return "invalid"
	}
}

type TicketSink interface {
	Create(goup, title, desc string) (tid.Tid, error)
	Update(tid tid.Tid, desc string) error
	Close(tid tid.Tid, desc string) error
	Group() string
}

type Metrics struct {
	tickets		*promp.CounterVec
	errors		promp.Counter
}
type Ticket struct {
	Default			TSys
	CloseResolved	bool
	metrics			Metrics
	MetrTicketsGend	*promp.CounterVec
	MetrErrors		promp.Counter
	sinks			[]TicketSink
	debug			bool
}

func New(cfg config.Ticket, dbg bool) *Ticket {

	t := &Ticket{

		debug:			dbg,
		Default:		NewTSys(cfg.Default),
		CloseResolved:	cfg.Resolved,
		metrics:		Metrics{
			tickets:	proma.NewCounterVec(
				promp.CounterOpts{
					Namespace:	"agate",
					Subsystem:	"ticket",
					Name:		"generated",
					Help:		"number of ticekts created",
				}, []string{
					"sys",
					"grp",
				}),
			errors: proma.NewCounter(
				promp.CounterOpts{
					Namespace:	"agate",
					Subsystem:	"ticket",
					Name:		"errors",
					Help:		"number of ticket errors",
				}),
		},
	}
	t.sinks = make([]TicketSink,len(tsysmap))
	t.sinks[TSysMock]	= mock.New(cfg.Sys.Mock,TSysMock.Int(),dbg)
	t.sinks[TSysGitlab] = gitlab.New(cfg.Sys.Gitlab,TSysGitlab.Int(),dbg)
	t.sinks[TSysHpsm]	= hpsm.New(cfg.Sys.Hpsm,TSysHpsm.Int(),dbg)

	if TSysMock > t.Default || t.Default >= TSysUnknown {
		t.unregister()
		panic(fmt.Sprintf("invalid default ticket sys: %d",t.Default))
	}

	return t
}

func (t *Ticket) NewTSysString(s string) TSys {
	if v, ok := tsysmap[s]; ok {
		return v
	} else {
		return t.Default
	}
}


func (t *Ticket) Del() {
	t.unregister()
}

func (t *Ticket) unregister() {
	if t != nil &&  t.metrics.errors != nil {
		promp.Unregister(t.metrics.errors);
		t.metrics.errors = nil
	}
	if t.metrics.tickets != nil  {
		promp.Unregister(t.metrics.tickets);
		t.metrics.tickets = nil
	}
}
func (t *Ticket) Sink(s TSys) TicketSink {

	if TSysMock <= s && s < TSysUnknown {
		return t.sinks[s]
	} else {
		return nil
	}
}

func (t *Ticket) Group(s TSys) string {
	if t.Sink(s) != nil {
		return t.Sink(s).Group()
	} else {
		return "invalid"
	}
}

func (t *Ticket) Create(s TSys, grp, title,desc string) tid.Tid {
	if t.Sink(s) != nil {
		stid, err := t.Sink(s).Create(grp,title,desc)
		if err == nil {
			return stid
		} else if t.debug  {
			panic(err)
		}
	}
	return nil
}

func (t *Ticket) Update(id tid.Tid, msg string) bool {
	if t.Sink(TSys(id.Sys())) != nil {
		err := t.Sink(TSys(id.Sys())).Update(id,msg)
		if err == nil {
			return true
		} else if t.debug  {
			panic(err)
		}
	}
	return false
}

func (t *Ticket) Close(id tid.Tid, msg string) bool {
	if t.Sink(TSys(id.Sys())) != nil {
		err := t.Sink(TSys(id.Sys())).Close(id,msg)
		if err == nil {
			return true
		} else if t.debug  {
			panic(err)
		}
	}
	return false
}

func (t *Ticket) Errorf(format string, args ...interface{}) error {
	t.MetrErrors.Inc()
	return fmt.Errorf(format,args...)
}
