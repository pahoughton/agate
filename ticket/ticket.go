/* 2018-12-27 (cc) <paul4hough@gmail.com>
   ticket management interface
*/
package ticket

import (
	"errors"

	"github.com/pahoughton/agate/config"
	"github.com/pahoughton/agate/ticket/gitlab"
	"github.com/pahoughton/agate/ticket/mock"

	proma "github.com/prometheus/client_golang/prometheus/promauto"
	promp "github.com/prometheus/client_golang/prometheus"
)
type Ticket struct {
	DefaultSys	string
	DefaultGrp	string
	Gitlab		*gitlab.Gitlab
	Mock		*mock.Mock
	TicketsGend	*promp.CounterVec
}

func New(c *config.Config) *Ticket {
	tck := &Ticket{
		DefaultSys: c.TicketDefaultSys,
		DefaultGrp:	c.TicketDefaultGrp,
		Gitlab:	gitlab.New(c.GitlabURL, c.GitlabToken, c.GitlabProject),
		Mock:	mock.New(c.MockURL, c.Debug),

		TicketsGend: proma.NewCounterVec(
			promp.CounterOpts{
				Namespace: "agate",
				Name:      "tickets_generated_total",
				Help:      "number of ticekts created",
			}, []string{
				"type",
				"dest",
			}),
	}

	tck.Gitlab.Debug = c.Debug

	return tck
}

func (t *Ticket) Create(
	tsys	string,
	tsub	string,
	title	string,
	desc	string ) (string, error) {

	var tid  string
	var err  error

	switch tsys {
	case "gitlab":
		tid, err = t.Gitlab.CreateIssue(tsub,title,desc)
	case "mock":
		tid, err = t.Mock.Create(title,desc)
	default:
		err = errors.New("unsupported ticket sys: " + tsys)
	}
	if err != nil {
		return "", err
	}

	return tid, nil
}

func (t *Ticket)AddComment(tsys string, tid string, cmt string ) error {

	switch tsys {
	case "gitlab":
		return t.Gitlab.AddComment(tid,cmt)
	case "mock":
		return t.Mock.AddComment(tid,cmt)
	default:
		return errors.New("unsupported ticket sys: "+tsys)
	}
}

func (t *Ticket)Close(tsys string, tid string ) error {

	switch tsys {
	case "gitlab":
		return t.Gitlab.Close(tid)
	case "mock":
		return t.Mock.Close(tid)
	default:
		return errors.New("unsupported ticket sys: "+tsys)
	}
}
