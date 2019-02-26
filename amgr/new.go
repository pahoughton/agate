/* 2018-12-25 (cc) <paul4hough@gmail.com>
   Prometheus AlertManager Alerts Body
*/

package amgr

import (
	"fmt"
	"sync"
	"time"

	proma "github.com/prometheus/client_golang/prometheus/promauto"
	promp "github.com/prometheus/client_golang/prometheus"

	"github.com/pahoughton/agate/config"
 	"github.com/pahoughton/agate/ticket"
	"github.com/pahoughton/agate/remed"
	"github.com/pahoughton/agate/db"
)

const (
	ATimeFmt = "2006-01-02T15:04:05.000000000-07:00"
)

type Metrics struct {
    groups	*promp.CounterVec
	alerts	*promp.CounterVec
	remed	*promp.GaugeVec
	errors	promp.Counter
}
type Fix struct {
	cnt		int32
	max		uint
	wg		sync.WaitGroup
	remed	*remed.Remed
}

type Amgr struct {
	debug	bool
	retry	time.Duration
	db		*db.DB
	qmgr	*Manager
	ticket	*ticket.Ticket
	fix		Fix
	metrics	Metrics
}

func New(c *config.Config,dataDir string,dbg bool) *Amgr {

	adb, err := db.New(dataDir, 0664, c.Global.DataAge,dbg);
	if err != nil {
		panic(err)
	}
	am := &Amgr{
		debug:		dbg,
		retry:		c.Global.Retry,
		db:			adb,
		qmgr:		NewManager(),
		ticket:		ticket.New(c.Ticket,dbg),
		fix:		Fix{
			max:		c.Global.Remed,
			remed:		remed.New(c.Global,dbg),
		},
		metrics: Metrics{
			groups: proma.NewCounterVec(
				promp.CounterOpts{
					Namespace: "agate",
					Subsystem:	"amgr",
					Name:      "groups_recvd",
					Help:      "number of alert groups received",
				},[]string{"resolve"}),
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
func (am *Amgr) Close() {
	am.ticket.Del()
	am.fix.remed.Close()
	am.db.Close()
	am.unregister()
}
func (am *Amgr) Errorf(format string, args ...interface{}) error {
	am.metrics.errors.Inc()
	return fmt.Errorf(format,args...)
}
func (am *Amgr) Error(err error) {
	am.metrics.errors.Inc()
	fmt.Println("ERROR: ",err.Error())
	if am.debug { panic(err); }
}
