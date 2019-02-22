/* 2019-01-07 (cc) <paul4hough@gmail.com>
   create remed instance
*/
package remed

import (
	"fmt"

	promp "github.com/prometheus/client_golang/prometheus"
	proma "github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/pahoughton/agate/config"
)

type Metrics struct {
	ansible		*promp.CounterVec
	scripts		*promp.CounterVec
	errors		promp.Counter
}
type Remed struct {
	debug			bool
	playbookDir		string
	scriptsDir		string
	metrics			Metrics
}

func New(c config.Global, dbg bool) *Remed {
	r := &Remed{
		debug:			dbg,
		playbookDir:	c.PlaybookDir,
		scriptsDir:		c.ScriptsDir,
		metrics:		Metrics{
			ansible: proma.NewCounterVec(
				promp.CounterOpts{
					Namespace:	"agate",
					Subsystem:	"remed",
					Name:		"ansible",
					Help:		"number of ansible playbook runs",
				}, []string{
					"role",
					"status",
				}),
			scripts: proma.NewCounterVec(
				promp.CounterOpts{
					Namespace: "agate",
					Subsystem: "remed",
					Name:      "script",
					Help:      "number of script runs",
				}, []string{
					"script",
					"status",
				}),
			errors: proma.NewCounter(
				promp.CounterOpts{
					Namespace: "agate",
					Subsystem: "remed",
					Name:      "errors",
					Help:      "number of errors",
				}),
		},
	}

	return r
}
func (r *Remed) Close() {
	r.unregister()
}

func (r *Remed) unregister() {
	if r != nil &&  r.metrics.ansible != nil {
		promp.Unregister(r.metrics.ansible);
		r.metrics.ansible = nil
	}
	if r.metrics.scripts != nil  {
		promp.Unregister(r.metrics.scripts);
		r.metrics.scripts = nil
	}
	if r.metrics.errors != nil  {
		promp.Unregister(r.metrics.errors);
		r.metrics.errors = nil
	}
}

func (r *Remed) Errorf(msg string,args ...interface{}) error {
	r.metrics.errors.Inc()
	return fmt.Errorf(msg,args...)
}
