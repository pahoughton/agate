/* 2019-01-07 (cc) <paul4hough@gmail.com>
   FIXME what is this for?
*/
package proc

import (
	"gitlab.com/pahoughton/agate/config"
	"gitlab.com/pahoughton/agate/ticket"

	promp "github.com/prometheus/client_golang/prometheus"
	proma "github.com/prometheus/client_golang/prometheus/promauto"
)

type Proc struct {
	Debug			bool
	Ticket			*ticket.Ticket
	PlaybookDir		string
	ScriptsDir		string
	AnsiblePlays	*promp.CounterVec
	ScriptsRun		*promp.CounterVec
}

func New(c *config.Config, t *ticket.Ticket) *Proc {
	p := &Proc{
		Debug:			c.Debug,
		Ticket:			t,
		PlaybookDir:	c.PlaybookDir,
		ScriptsDir:		c.ScriptsDir,
		AnsiblePlays: proma.NewCounterVec(
			promp.CounterOpts{
				Namespace: "agate",
				Name:      "ansible_plays_total",
				Help:      "number of ansible playbook runs",
			}, []string{
				"role",
				"status",
			}),
		ScriptsRun: proma.NewCounterVec(
			promp.CounterOpts{
				Namespace: "agate",
				Name:      "script_runs_total",
				Help:      "number of script runs",
			}, []string{
				"script",
				"status",
			}),
	}
	return p
}
