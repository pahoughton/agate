/* 2018-12-27 (cc) <paul4hough@gmail.com>
   ticket management interface
*/
package ticket

import (
	"errors"

	"github.com/pahoughton/agate/config"
	"github.com/pahoughton/agate/model"
	"github.com/pahoughton/agate/ticket/gitlab"
	"github.com/pahoughton/agate/ticket/hpsm"
	"github.com/pahoughton/agate/ticket/mock"

	proma "github.com/prometheus/client_golang/prometheus/promauto"
	promp "github.com/prometheus/client_golang/prometheus"
)
type Ticket struct {
	Debug		bool
	DefaultSys	string
	DefaultGrp	string
	Gitlab		*gitlab.Gitlab
	Mock		*mock.Mock
	HPSM		*hpsm.HPSM
	TicketsGend	*promp.CounterVec
}

func New(c *config.Config, dbg bool) *Ticket {

	t := &Ticket{
		Debug:		dbg,
		DefaultSys: c.TicketDefaultSys,
		DefaultGrp:	c.TicketDefaultGrp,
		Gitlab:		gitlab.New(c.GitlabURL,c.GitlabToken,dbg),
		HPSM:		hpsm.New(c.HpsmURL,c.HpsmUser,c.HpsmPass,dbg),
		Mock:		mock.New(c.MockURL,dbg),

		TicketsGend: proma.NewCounterVec(
			promp.CounterOpts{
				Namespace: "agate",
				Name:      "tickets_generated_total",
				Help:      "number of ticekts created",
			}, []string{
				"sys",
				"grp",
			}),
	}

	return t
}

func (t *Ticket) Sys(a model.Alert) string {
	sys := string(a.Annotations["ticket"])
	if len(sys) < 1 {
		sys = t.DefaultSys
	}
	return sys
}

func (t *Ticket) Grp(a model.Alert) string {
	grp := string(a.Annotations["ticket_group"])
	if len(grp) < 1 {
		grp = t.DefaultGrp
	}
	return grp
}

func (t *Ticket) Create(a model.Alert) (string, error) {

	var tid  string
	var err  error

	switch t.Sys(a) {
	case "hpsm":
		tid, err = t.HPSM.Create(t.Grp(a),a)
	case "gitlab":
		tid, err = t.Gitlab.Create(t.Grp(a),a)
	case "mock":
		tid, err = t.Mock.Create(a)
	default:
		err = errors.New("unsupported ticket sys: " + t.Sys(a))
	}

	return tid,err
}

func (t *Ticket)AddComment(a model.Alert, tid, cmt string ) error {

	switch t.Sys(a) {
	case "hpsm":
		return t.HPSM.AddComment(tid,cmt)
	case "gitlab":
		return t.Gitlab.AddComment(tid,cmt)
	case "mock":
		return t.Mock.AddComment(tid,cmt)
	default:
		return errors.New("unsupported ticket sys: "+t.Sys(a))
	}
}

func (t *Ticket)Close(a model.Alert,tid string) error {

	switch t.Sys(a) {
	case "hpsm":
		return t.HPSM.Close(tid)
	case "gitlab":
		return t.Gitlab.Close(tid)
	case "mock":
		return t.Mock.Close(tid)
	default:
		return errors.New("unsupported ticket sys: "+t.Sys(a))
	}
}
