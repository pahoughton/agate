/* 2019-02-19 (cc) <paul4hough@gmail.com>
   FIXME what is this for?
*/
package db

import (
	"os"
	"path"
	"testing"
	"github.com/stretchr/testify/require"
)

func testit(t *testing.T, check func(t *testing.T,db *DB)) {

	os.Remove(path.Join("testdata",dbFn))
	db, err := New("testdata",0664,5,true)
	require.Nil(t,err)
	require.NotNil(t,db)
	check(t,db)
	db.Del()
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
		retCode := m.Run()
		os.Exit(retCode)
	} else {
		panic(err)
	}
}
