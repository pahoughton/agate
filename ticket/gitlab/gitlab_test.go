/* 2019-02-17 (cc) <paul4hough@gmail.com>
   validate gitlab ticket (issue) interface
*/
package gitlab
import (
	"testing"
	"net"
	"net/http/httptest"

	"github.com/stretchr/testify/assert"

	"github.com/pahoughton/agate/config"
	"github.com/pahoughton/agate/ticket/tid"

)

const (
	Token	= "secret-sauce"
)
func TestNew(t *testing.T) {

	gl := New(
		config.TSysGitlab{
			Url:	"http://gitlab.com",
			Token:	Token,
			Group:	"test",
		},2,false)
	assert.NotNil(t,gl.c)
	assert.False(t,gl.debug)
}
func TestGroup(t *testing.T) {

	exp := "paul/test"
	gl := New(
		config.TSysGitlab{
			Url:	"http://gitlab.com",
			Token:	Token,
			Group:	exp,
		},2,false)
	assert.NotNil(t,gl.c)
	assert.Equal(t,exp,gl.Group())
}


func TestCreate(t *testing.T) {
	mock := NewMockServer()
	ms := httptest.NewServer(mock)
	defer ms.Close()

	expGrp := "paul/test"
	gl := New(
		config.TSysGitlab{
			Url:	ms.URL,
			Token:	Token,
			Group:	expGrp,
		},2,false)
	tid, err := gl.Create(expGrp,"broken stuff","details details")
	assert.NotNil(t,tid)
	assert.Nil(t,err)
	assert.Equal(t,expGrp + ":1",tid.String())
}

func TestCreateNetError(t *testing.T) {
	gl := New(
		config.TSysGitlab{
			Url:	"http://localhost:31001/ticket",
			Token:	Token,
			Group:	"test",
		},2,false)
	tid, err := gl.Create("storage","disk full","disk is full")
	assert.Nil(t,tid)
	assert.NotNil(t,err)
	_, ok := err.(net.Error);
	assert.True(t,ok)
}


func TestUpdate(t *testing.T) {
	mock := NewMockServer()
	ms := httptest.NewServer(mock)
	defer ms.Close()

	gl := New(
		config.TSysGitlab{
			Url:	ms.URL,
			Token:	Token,
			Group:	"test",
		},2,false)
	err := gl.Update(tid.NewString(2,"paul/test:14"),"comment")
	assert.Nil(t,err)
}

func TestUpdateNetError(t *testing.T) {
	gl := New(
		config.TSysGitlab{
			Url:	"http://localhost:31001/ticket",
			Token:	Token,
			Group:	"test",
		},2,false)
	err := gl.Update(tid.NewString(2,"prj:12"),"disk still full")
	assert.NotNil(t,err)
	_, ok := err.(net.Error);
	assert.True(t,ok)
}

func TestClose(t *testing.T) {
	mock := NewMockServer()
	ms := httptest.NewServer(mock)
	defer ms.Close()

	gl := New(
		config.TSysGitlab{
			Url:	ms.URL,
			Token:	Token,
			Group:	"test",
		},2,false)
	err := gl.Close(tid.NewString(2,"paul/test:14"),"close")
	assert.Nil(t,err)
}

func TestCloseNetError(t *testing.T) {
	gl := New(
		config.TSysGitlab{
			Url:	"http://localhost:31001/ticket",
			Token:	Token,
			Group:	"test",
		},2,false)
	err := gl.Close(tid.NewString(2,"prj:12"),"fixed")
	assert.NotNil(t,err)
	_, ok := err.(net.Error);
	assert.True(t,ok)
}
