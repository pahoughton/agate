/* 2018-12-25 (cc) <paul4hough@gmail.com>
   Prometheus AlertManager Alerts Body
*/

package amgr

import (
	"fmt"
	"time"
	proma "github.com/prometheus/client_golang/prometheus/promauto"
	promp "github.com/prometheus/client_golang/prometheus"

	"github.com/pahoughton/agate/config"
 	"github.com/pahoughton/agate/notify"
	"github.com/pahoughton/agate/remed"
)

type Metrics struct {
    groups	*promp.CounterVec
	alerts	*promp.CounterVec
	remed	*promp.GaugeVec
	errors	promp.Counter
}

type Amgr struct {
	debug	bool
	notify	*notify.Notify
	remed	*remed.Remed
	retry	time.Duration
	metrics	Metrics
}

func New(c *config.Config,dataDir string,dbg bool) *Amgr {

	notify := notify.New(c.Notify,dataDir,dbg)
	am := &Amgr{
		debug:		dbg,
		retry:		c.Global.Retry,
		notify:		notify,
		remed:		remed.New(c.Remed,dataDir,notify,dbg),
		metrics: Metrics{
			groups: proma.NewCounterVec(
				promp.CounterOpts{
					Namespace: "agate",
					Subsystem:	"amgr",
					Name:      "groups_recvd",
					Help:      "number of alert groups received",
				},[]string{"sys","grp","resolve"}),
			alerts: proma.NewCounterVec(
				promp.CounterOpts{
					Namespace: "agate",
					Subsystem:	"amgr",
					Name:      "alerts_recvd",
					Help:      "number of alerts received",
				}, []string{
					"name",
					"node",
					"resolve",
				}),
			remed: proma.NewGaugeVec(
				promp.GaugeOpts{
					Namespace: "agate",
					Subsystem:	"amgr",
					Name:      "remed_active",
					Help:      "number of running remediation attempts",
				}, []string{"alert"}),
			errors: proma.NewCounter(
				promp.CounterOpts{
					Namespace: "agate",
					Subsystem:	"amgr",
					Name:      "errors",
					Help:      "number of amgr errors",
				}),
		},
	}
	return am
}

func (am *Amgr) Del() {
	am.notify.Del()
	am.remed.Del()
	am.unregister()
}

func (am *Amgr) unregister() {
	if am.metrics.groups != nil {
		promp.Unregister(am.metrics.groups)
		am.metrics.groups = nil
	}
	if am.metrics.alerts != nil {
		promp.Unregister(am.metrics.alerts)
		am.metrics.alerts = nil
	}
	if am.metrics.remed != nil {
		promp.Unregister(am.metrics.remed)
		am.metrics.remed = nil
	}
	if am.metrics.errors != nil {
		promp.Unregister(am.metrics.errors)
		am.metrics.errors = nil
	}
}
func (am *Amgr) errorf(format string, args ...interface{}) error {
	am.metrics.errors.Inc()
	return fmt.Errorf(format,args...)
}
func (am *Amgr) error(err error) {
	am.metrics.errors.Inc()
	fmt.Println("ERROR: ",err.Error())
	if am.debug { panic(err); }
}
