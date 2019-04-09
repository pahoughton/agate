/* 2019-02-17 (cc) <paul4hough@gmail.com>
   validate mock ticket interface
*/
package mock

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pahoughton/agate/config"
	"github.com/pahoughton/agate/notify/nid"
)

func TestNew(t *testing.T) {
	m := New(config.NSysMock{Url: "http://localhost:5001/ticket"},0,false)
	assert.NotNil(t,m)
}
func TestGroup(t *testing.T) {
	m := New(config.NSysMock{Url: "http://localhost:5001/ticket"},0,false)
	assert.Equal(t,"",m.Group())
}

func TestCreate(t *testing.T) {
	mock := &MockServer{}
	msrv := httptest.NewServer(mock)
	defer msrv.Close()

	m := New(config.NSysMock{Url: msrv.URL},0,false)
	assert.NotNil(t,m)
	nid, err := m.Create("storage","disk full","disk is full")
	assert.Nil(t,err)
	assert.NotNil(t,nid)
	assert.Equal(t,"1",nid.Id())
}
func TestCreateSysError(t *testing.T) {
	respJson := `garbage`

	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, respJson)
			}))
	defer ts.Close()

	m := New(config.NSysMock{Url: ts.URL},0,false)
	assert.Panics(t, func() {
		m.Create("storage","disk full","disk is full")
	}, "create bad resp should panic")
}

func TestCreateNetError(t *testing.T) {
	m := New(config.NSysMock{Url: "http://localhost:31001/ticket"},0,false)
	assert.NotNil(t,m)
	nid, err := m.Create("storage","disk full","disk is full")
	assert.Nil(t,nid)
	assert.NotNil(t,err)
	_, ok := err.(net.Error);
	assert.True(t,ok)
}

func TestUpdate(t *testing.T) {
	mock := &MockServer{}
	msrv := httptest.NewServer(mock)
	defer msrv.Close()

	m := New(config.NSysMock{Url: msrv.URL},0,false)
	assert.NotNil(t,m)
	err := m.Update(nid.NewString(1,"12"),"disk still full")
	assert.Nil(t,err)
}

func TestUpdateNetError(t *testing.T) {
	m := New(config.NSysMock{Url: "http://localhost:31001/ticket"},0,false)
	assert.NotNil(t,m)
	err := m.Update(nid.NewString(0,"12"),"disk still full")
	assert.NotNil(t,err)
	_, ok := err.(net.Error);
	assert.True(t,ok)
}

func TestClose(t *testing.T) {
	mock := &MockServer{}
	msrv := httptest.NewServer(mock)
	defer msrv.Close()

	m := New(config.NSysMock{Url: msrv.URL},0,false)
	assert.NotNil(t,m)
	err := m.Close(nid.NewString(0,"12"),"fixed")
	assert.Nil(t,err)
}

func TestCloseNetError(t *testing.T) {
	m := New(config.NSysMock{Url: "http://localhost:31001/ticket"},0,false)
	assert.NotNil(t,m)
	err := m.Close(nid.NewString(0,"12"),"fixed")
	assert.NotNil(t,err)
	_, ok := err.(net.Error);
	assert.True(t,ok)
}
