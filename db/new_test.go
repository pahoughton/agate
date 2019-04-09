/* 2019-02-19 (cc) <paul4hough@gmail.com>
   FIXME what is this for?
*/
package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)
func TestNew(t *testing.T) {

	db, err := New("testdata",0664,15,false)
	assert.Nil(t,err)
	assert.NotNil(t,db)
	edb, err := New("testdata",0664,15,false)
	assert.NotNil(t,err)
	assert.Nil(t,edb)
	db.Del()
	db, err = New("testdata",0664,15,false)
	assert.Nil(t,err)
	assert.NotNil(t,db)
	db.Del()
}
