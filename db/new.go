/* 2018-12-31 (cc) <paul4hough@gmail.com>
   agate alert db stores ticket id's for alerts so that tickets for resolved
   alerts can be updated.

   Each date will have it's own buckets to provide for deleting
   unresolved alerts that are older than 'MaxAge' days

*/
package db

import (
	"fmt"
	"os"
	"path"
	"time"

	promp "github.com/prometheus/client_golang/prometheus"
	"github.com/boltdb/bolt"
)

const (
	dbFn			= "agate.bolt"
)

type Metrics struct {
	agqueue		promp.Gauge
	tickets		*promp.GaugeVec
	errors		promp.Counter
}
type DB struct {
	debug	bool
	maxDays int
	db		*bolt.DB
	metrics	*Metrics
}

func New(dir string, mode os.FileMode, maxDays uint,debug bool) (*DB, error) {

	fn := path.Join(dir,dbFn)
	_, err := os.Stat(fn)
	isnew := err != nil || os.IsNotExist(err)
	opts := &bolt.Options{Timeout: 50 * time.Millisecond}
	bdb, err := bolt.Open(fn,mode,opts)
	if err != nil {
		return nil, err
	}
	if isnew {
		err := bdb.Update(func(tx *bolt.Tx) error {
			if _, err := tx.CreateBucket(agroupBucket()); err != nil {
				return err
			}
			if _, err := tx.CreateBucket(alertsBucket()); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	db := &DB{
		debug: debug,
		db: bdb,
		maxDays: int(maxDays),
		metrics: &Metrics{
			agqueue: promp.NewGauge(
				promp.GaugeOpts{
					Namespace: "agate",
					Subsystem: "db",
					Name:      "agroup_queue_size",
					Help:      "number of records in agroup bucket",
				}),
			tickets: promp.NewGaugeVec(
				promp.GaugeOpts{
					Namespace: "agate",
					Subsystem: "db",
					Name:      "open_alerts_size",
					Help:      "number of records in agroup bucket",
				},[]string{ "date" }),
			errors: promp.NewCounter(
				promp.CounterOpts{
					Namespace: "agate",
					Subsystem: "db",
					Name:      "errors_total",
					Help:      "number of records in agroup bucket",
				}),
		},
	}

	promp.MustRegister(
		db.metrics.agqueue,
		db.metrics.tickets,
		db.metrics.errors)

	db.AlertCleanBuckets()
	// reclean alert buckets every 24 hours
	go func() {
		for _ = range time.NewTicker(time.Hour * 24).C {
			fmt.Println("INFO cleaning buckets before")
			db.AlertCleanBuckets()
		}
	}()

	return db, nil
}
func (db *DB) Close() {
	db.unregister()
	if db.db != nil { db.db.Close(); db.db = nil }
}

func (db *DB) unregister() { // for testing
	if db != nil && db.metrics != nil && db.metrics.agqueue != nil {
		promp.Unregister(db.metrics.agqueue);
		db.metrics.agqueue = nil
	}
	if db.metrics.tickets != nil  {
		promp.Unregister(db.metrics.tickets);
		db.metrics.tickets = nil
	}
	if db.metrics.errors != nil  {
		promp.Unregister(db.metrics.errors);
		db.metrics.errors = nil
	}
}
