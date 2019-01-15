/* 2018-12-27 (cc) <paul4hough@gmail.com>
   ticket management interface
*/
package ticket

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/pahoughton/agate/db"
	"github.com/pahoughton/agate/config"
	"github.com/pahoughton/agate/ticket/gitlab"
	"github.com/pahoughton/agate/ticket/mock"

	proma "github.com/prometheus/client_golang/prometheus/promauto"
	promp "github.com/prometheus/client_golang/prometheus"
)
type Ticket struct {
	DefaultSys	string
	DefaultGrp	string
	Adb			*db.AlertDB
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

	var err error
	tck.Adb, err = db.Open(path.Join(c.BaseDir, "data"), 0664, c.MaxDays);
	if err != nil {
		fmt.Println("FATAL: open db - ",err.Error())
		os.Exit(1)
	}

	return tck
}

func (t *Ticket) Create(
	tsys	string,
	tsub	string,
	aKey	string,
	title	string,
	desc	string ) (string, error) {

	var tid  string
	var err  error

	if len(tsys) < 1 {
		tsys = t.DefaultSys
	}
	if len(tsub) < 1 {
		tsub = t.DefaultGrp
	}
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

	return tid, t.Adb.AddTicket(aKey,tid)
}

func (t *Ticket)AddTidComment(tsys string, tid string, cmt string ) error {

	if len(tsys) < 1 {
		tsys = t.DefaultSys
	}
	switch tsys {
	case "gitlab":
		return t.Gitlab.AddComment(tid,cmt)
	case "mock":
		return t.Mock.AddComment(tid,cmt)
	default:
		return errors.New("unsupported ticket sys: "+tsys)
	}
}

func (t *Ticket)AddKeyComment(tsys string, aKey string, cmt string ) error {

	if len(tsys) < 1 {
		tsys = t.DefaultSys
	}
	tid, err := t.Adb.GetTicket(aKey)
	if err != nil {
		return err
	}
	switch tsys {
	case "gitlab":
		return t.Gitlab.AddComment(tid,cmt)
	case "mock":
		return t.Mock.AddComment(tid,cmt)
	default:
		return errors.New("unsupported ticket sys: "+tsys)
	}
}

func (t *Ticket)Close(tsys string, aKey string ) error {

	if len(tsys) < 1 {
		tsys = t.DefaultSys
	}
	tid, err := t.Adb.GetTicket(aKey)
	if err != nil {
		return err
	}
	switch tsys {
	case "gitlab":
		return t.Gitlab.Close(tid)
	case "mock":
		return t.Mock.Close(tid)
	default:
		return errors.New("unsupported ticket sys: "+tsys)
	}
}

func (t *Ticket)Delete(aKey string ) error {
	return t.Adb.DelTicket(aKey)
}
