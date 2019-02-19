/* 2019-02-17 (cc) <paul4hough@gmail.com>
   validate gitlab ticket (issue) interface
*/
package gitlab
import (
	"fmt"
	"testing"
	"net"
	"net/http"
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
		},
		false)
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
		},
		false)
	assert.NotNil(t,gl.c)
	assert.Equal(t,exp,gl.Group())
}


func TestCreate(t *testing.T) {
	respJson := `{"id":1, "iid":14, "title" : "Title of issue",
"description": "This is description of an issue",
"author" : {"id" : 1, "name": "snehal"}, "assignees":[{"id":1}]}`

	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, respJson)
	}))
	defer ts.Close()

	expGrp := "paul/test"
	gl := New(
		config.TSysGitlab{
			Url:	ts.URL,
			Token:	Token,
			Group:	expGrp,
		},
		false)
	tid, err := gl.Create(expGrp,"broken stuff","details details")
	assert.NotNil(t,tid)
	assert.Nil(t,err)
	assert.Equal(t,expGrp + ":14",tid.String())
}

func TestCreateNetError(t *testing.T) {
	gl := New(
		config.TSysGitlab{
			Url:	"http://localhost:31001/ticket",
			Token:	Token,
			Group:	"test",
		},
		false)
	tid, err := gl.Create("storage","disk full","disk is full")
	assert.Nil(t,tid)
	assert.NotNil(t,err)
	_, ok := err.(net.Error);
	assert.True(t,ok)
}


func TestUpdate(t *testing.T) {
	respJson := `{"id":1, "iid":14, "title" : "Title of issue",
"description": "This is description of an issue",
"author" : {"id" : 1, "name": "snehal"}, "assignees":[{"id":1}]}`

	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, respJson)
	}))
	defer ts.Close()

	gl := New(
		config.TSysGitlab{
			Url:	ts.URL,
			Token:	Token,
			Group:	"test",
		},
		false)
	err := gl.Update(tid.NewString("paul/test:14"),"comment")
	assert.Nil(t,err)
}

func TestUpdateNetError(t *testing.T) {
	gl := New(
		config.TSysGitlab{
			Url:	"http://localhost:31001/ticket",
			Token:	Token,
			Group:	"test",
		},
		false)
	err := gl.Update(tid.NewString("prj:12"),"disk still full")
	assert.NotNil(t,err)
	_, ok := err.(net.Error);
	assert.True(t,ok)
}

func TestClose(t *testing.T) {
	respJson := `{"id":1, "iid":14, "title" : "Title of issue",
"description": "This is description of an issue",
"author" : {"id" : 1, "name": "snehal"}, "assignees":[{"id":1}]}`

	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, respJson)
	}))
	defer ts.Close()

	gl := New(
		config.TSysGitlab{
			Url:	ts.URL,
			Token:	Token,
			Group:	"test",
		},
		false)
	err := gl.Close(tid.NewString("paul/test:14"),"close")
	assert.Nil(t,err)
}

func TestCloseNetError(t *testing.T) {
	gl := New(
		config.TSysGitlab{
			Url:	"http://localhost:31001/ticket",
			Token:	Token,
			Group:	"test",
		},
		false)
	err := gl.Close(tid.NewString("prj:12"),"fixed")
	assert.NotNil(t,err)
	_, ok := err.(net.Error);
	assert.True(t,ok)
}
