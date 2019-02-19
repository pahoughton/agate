/* 2018-12-27 (cc) <paul4hough@gmail.com>
   ticket management interface
*/
package ticket

import (
	"fmt"

	"github.com/pahoughton/agate/config"
	"github.com/pahoughton/agate/amgr/alert"
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

func (t TSys) String() string {

	names := []string{"unk","gitlab","hpsm","mock"}
	if TSysMock <= t && t <= TSysUnknown {
		return names[t]
	} else {
		return "invalid"
	}
}

type TicketSink interface {
	Create(goup, title, desc string) (*tid.Tid, error)
	Update(tid *tid.Tid, desc string) error
	Close(tid *tid.Tid, desc string) error
	Group() string
}

type AlertTid interface {
	Tid() []byte
}
type Ticket struct {
	sinks			[]TicketSink
	Debug			bool
	Default			TSys
	CloseResolved	bool
	MetrTicketsGend	*promp.CounterVec
	MetrErrors		promp.Counter
}

func New(cfg config.Ticket, dbg bool) *Ticket {

	t := &Ticket{

		Debug:			dbg,
		Default:		NewTSys(cfg.Default),
		CloseResolved:	true,
		MetrTicketsGend: proma.NewCounterVec(
			promp.CounterOpts{
				Namespace: "agate",
				Name:      "tickets_generated_total",
				Help:      "number of ticekts created",
			}, []string{
				"sys",
				"grp",
			}),
		MetrErrors: proma.NewCounter(
			promp.CounterOpts{
				Namespace: "agate",
				Name:      "ticket_errors_total",
				Help:      "number of ticket errors",
			}),
	}
	t.sinks = make([]TicketSink,len(tsysmap))
	t.sinks[TSysMock]	= mock.New(cfg.Sys.Mock,dbg)
	t.sinks[TSysGitlab] = gitlab.New(cfg.Sys.Gitlab,dbg)
	t.sinks[TSysHpsm]	= hpsm.New(cfg.Sys.Hpsm,dbg)

	if TSysMock > t.Default || t.Default >= TSysUnknown {
		promp.Unregister(t.MetrTicketsGend)
		promp.Unregister(t.MetrErrors)
		panic(fmt.Sprintf("invalid default ticket sys: %d",t.Default))
	}
	return t
}

func (t *Ticket) Sink(s TSys) TicketSink {

	if TSysUnknown < s || s >= TSysMock {
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

func (t *Ticket) Errorf(format string, args ...interface{}) error {
	t.MetrErrors.Inc()
	return fmt.Errorf(format,args...)
}

func (t *Ticket) AlertTSys(a alert.Alert) TSys {

	tsys := t.Default
	lsys, ok := a.Labels["ticket_sys"]
	if ok {
		if tmp, ok := tsysmap[string(lsys)]; ok {
			tsys = tmp
		} else {
			t.Errorf("alert unknown tsys %v",lsys)
		}
	}
	return tsys
}

func (t *Ticket) AgroupTSys(agrp alert.AlertGroup) TSys {

	if v, ok := agrp.ComLabels["ticket_sys"]; ok {
		sys, ok := tsysmap[string(v)]
		if ok {
			return sys
		} else {
			t.Errorf("agroup invalid ticket_sys: %s",string(v))
			return t.Default
		}
	} else {
		// majority rule
		majName := t.Default.String()
		majCount := 0
		agtmap :=  make(map[string]int,len(agrp.Alerts))
		for _, a := range agrp.Alerts {
			sname := t.Default.String()
			if v, ok := a.Labels["ticket_sys"]; ok {
				if _, ok := tsysmap[string(v)]; ok {
					sname = string(v)
				} else {
					t.Errorf("alert invalid ticket_sys: %s",string(v))
				}
			}
			agtmap[sname] += 1

			if agtmap[sname] > majCount {
				majCount = agtmap[sname]
				majName = sname
			}
		}
		return tsysmap[majName]
	}
}


func (t *Ticket) AlertTGrp(a alert.Alert) string {
	var tgrp string
	lgrp, ok := a.Labels["ticket_grp"]
	if ok {
		tgrp = string(lgrp)
	} else {
		tgrp = t.Group(t.AlertTSys(a))
	}
	return tgrp
}

func (t *Ticket) AgroupTGrp(ag alert.AlertGroup) string {

	if v, ok := ag.ComLabels["ticket_grp"]; ok {
		return string(v)
	} else {
		// majority rule
		defgrp := t.Group(t.AgroupTSys(ag))
		majName := defgrp
		majCount := 0
		agtmap :=  make(map[string]int,len(ag.Alerts))
		for _, a := range ag.Alerts {
			gname := defgrp
			if v, ok := a.Labels["ticket_grp"]; ok {
				gname = string(v)
			}
			agtmap[gname] += 1

			if agtmap[gname] > majCount {
				majCount = agtmap[gname]
				majName = gname
			}
		}
		return majName
	}
}
