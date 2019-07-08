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

	proma "github.com/prometheus/client_golang/prometheus/promauto"
	promp "github.com/prometheus/client_golang/prometheus"
	"github.com/boltdb/bolt"
)

const (
	DbFn			= "agate-2.bolt"
)
var DbPrevFn = []string{"agate.bolt", "agate-1.bolt"}

type Metrics struct {
	agqueue		*promp.GaugeVec
	dbucket		*promp.GaugeVec
	errors		promp.Counter
}
type DB struct {
	debug	bool
	maxDays int
	db		*bolt.DB
	metrics	*Metrics
}

func New(dir string, mode os.FileMode, maxDays uint,debug bool) (*DB, error) {

	for _, pfn := range DbPrevFn {
		os.Remove(path.Join(dir,pfn))
	}
	fn := path.Join(dir,DbFn)
	// _, err := os.Stat(fn)
	// isnew := err != nil || os.IsNotExist(err)
	opts := &bolt.Options{Timeout: 50 * time.Millisecond}
	bdb, err := bolt.Open(fn,mode,opts)
	if err != nil {
		return nil, err
	}
	db := &DB{
		debug: debug,
		db: bdb,
		maxDays: int(maxDays),
		metrics: &Metrics{
			agqueue: proma.NewGaugeVec(
				promp.GaugeOpts{
					Namespace: "agate",
					Subsystem: "db",
					Name:      "agroup_queue_size",
					Help:      "number of records in agroup bucket",
				},[]string{"sys"}),
			dbucket: proma.NewGaugeVec(
				promp.GaugeOpts{
					Namespace: "agate",
					Subsystem: "db",
					Name:      "bucket_size",
					Help:      "number of records in agroup bucket",
				},[]string{ "date","bucket" }),
			errors: proma.NewCounter(
				promp.CounterOpts{
					Namespace: "agate",
					Subsystem: "db",
					Name:      "errors_total",
					Help:      "number of records in agroup bucket",
				}),
		},
	}

	db.CleanDateBuckets()
	go func() {
		for _ = range time.NewTicker(time.Hour * 24).C {
			fmt.Println("INFO cleaning buckets before")
			db.CleanDateBuckets()
		}
	}()

	return db, nil
}
func (db *DB) Del() {
	db.unregister()
	if db.db != nil { db.db.Close(); db.db = nil }
}

func (db *DB) unregister() { // for testing
	if db != nil && db.metrics != nil && db.metrics.agqueue != nil {
		promp.Unregister(db.metrics.agqueue);
		db.metrics.agqueue = nil
	}
	if db.metrics.dbucket != nil  {
		promp.Unregister(db.metrics.dbucket);
		db.metrics.dbucket = nil
	}
	if db.metrics.errors != nil  {
		promp.Unregister(db.metrics.errors);
		db.metrics.errors = nil
	}
}
