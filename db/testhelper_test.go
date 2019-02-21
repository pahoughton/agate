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

func benchit(b *testing.B, check func(b *testing.B,db *DB)) {

	os.Remove(path.Join("testdata",dbFn))
	db, err := New("testdata",0664,5,true)
	require.Nil(b,err)
	require.NotNil(b,db)
	check(b,db)
	db.Close()
}
func testit(t *testing.T, check func(t *testing.T,db *DB)) {

	os.Remove(path.Join("testdata",dbFn))
	db, err := New("testdata",0664,5,true)
	require.Nil(t,err)
	require.NotNil(t,db)
	check(t,db)
	db.Close()
}
