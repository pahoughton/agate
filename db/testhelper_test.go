/* 2019-02-19 (cc) <paul4hough@gmail.com>
   FIXME what is this for?
*/
package db

import (
	"os"
	"path"
	"testing"
	"github.com/stretchr/testify/require"
	promp "github.com/prometheus/client_golang/prometheus"
)

func testit(t *testing.T, check func(t *testing.T,db *DB)) {

	os.Remove(path.Join("testdata",dbFn))
	db, err := New("testdata",0664,5,true)
	require.Nil(t,err)
	require.NotNil(t,db)
	check(t,db)
	db.Close()
}
// global db for benchmark tests
var (
	benchdb *DB
)
func benchit(b *testing.B, check func(b *testing.B,db *DB)) {
	check(b,benchdb)
}
func TestMain(m *testing.M) {
	if db, err := New("testdata/bench",0664,5,true); err == nil {
		benchdb = db
		db.unregister()
		db.metrics = &Metrics{
			agqueue: promp.NewGauge(
				promp.GaugeOpts{
					Name: "bench_qsize",
					Help: "none",
				}),
			tickets: promp.NewGaugeVec(
				promp.GaugeOpts{
					Name: "bench_alerts",
					Help: "none",
				},[]string{ "date" }),
			errors: promp.NewCounter(
				promp.CounterOpts{
					Name: "bench_errors",
					Help: "none",
				}),
		}
		retCode := m.Run()
		db.Close()
		os.Exit(retCode)
	} else {
		panic(err)
	}
}
