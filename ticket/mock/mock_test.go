/* 2019-02-17 (cc) <paul4hough@gmail.com>
   validate mock ticket interface
*/
package mock

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pahoughton/agate/config"
	"github.com/pahoughton/agate/ticket/tid"
)

func TestNew(t *testing.T) {
	m := New(config.TSysMock{Url: "http://localhost:5001/ticket"},0,false)
	assert.NotNil(t,m)
}
func TestGroup(t *testing.T) {
	m := New(config.TSysMock{Url: "http://localhost:5001/ticket"},0,false)
	assert.Equal(t,"",m.Group())
}

func TestCreate(t *testing.T) {
	respJson := `{"id":"23"}`

	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				cont, err := ioutil.ReadAll(r.Body);
				if err != nil {
					t.Error("readall")
					return
				}
				var dat map[string]string
				if err := json.Unmarshal(cont, &dat); err != nil {
					t.Error("json.Unmarshal")
					return
				}
				fmt.Fprintln(w, respJson)
			}))
	defer ts.Close()

	m := New(config.TSysMock{Url: ts.URL},0,false)
	assert.NotNil(t,m)
	tid, err := m.Create("storage","disk full","disk is full")
	assert.Nil(t,err)
	assert.NotNil(t,tid)
	assert.Equal(t,"23",tid.String())
}
func TestCreateSysError(t *testing.T) {
	respJson := `garbage`

	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				cont, err := ioutil.ReadAll(r.Body);
				if err != nil {
					t.Error("readall")
					return
				}
				var dat map[string]string
				if err := json.Unmarshal(cont, &dat); err != nil {
					t.Error("json.Unmarshal")
					return
				}
				fmt.Fprintln(w, respJson)
			}))
	defer ts.Close()

	m := New(config.TSysMock{Url: ts.URL},0,false)
	assert.Panics(t, func() {
		m.Create("storage","disk full","disk is full")
	}, "create bad resp should panic")
}

func TestCreateNetError(t *testing.T) {
	m := New(config.TSysMock{Url: "http://localhost:31001/ticket"},0,false)
	assert.NotNil(t,m)
	tid, err := m.Create("storage","disk full","disk is full")
	assert.Nil(t,tid)
	assert.NotNil(t,err)
	_, ok := err.(net.Error);
	assert.True(t,ok)
}

func TestUpdate(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				cont, err := ioutil.ReadAll(r.Body);
				if err != nil {
					t.Error("readall")
					return
				}
				var dat map[string]string
				if err := json.Unmarshal(cont, &dat); err != nil {
					t.Error("json.Unmarshal")
					return
				}
			}))
	defer ts.Close()

	m := New(config.TSysMock{Url: ts.URL},0,false)
	assert.NotNil(t,m)
	err := m.Update(tid.NewString(1,"12"),"disk still full")
	assert.Nil(t,err)
}

func TestUpdateNetError(t *testing.T) {
	m := New(config.TSysMock{Url: "http://localhost:31001/ticket"},0,false)
	assert.NotNil(t,m)
	err := m.Update(tid.NewString(0,"12"),"disk still full")
	assert.NotNil(t,err)
	_, ok := err.(net.Error);
	assert.True(t,ok)
}

func TestClose(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				cont, err := ioutil.ReadAll(r.Body);
				if err != nil {
					t.Error("readall")
					return
				}
				var dat map[string]string
				if err := json.Unmarshal(cont, &dat); err != nil {
					t.Error("json.Unmarshal")
					return
				}
			}))
	defer ts.Close()

	m := New(config.TSysMock{Url: ts.URL},0,false)
	assert.NotNil(t,m)
	err := m.Close(tid.NewString(0,"12"),"fixed")
	assert.Nil(t,err)
}

func TestCloseNetError(t *testing.T) {
	m := New(config.TSysMock{Url: "http://localhost:31001/ticket"},0,false)
	assert.NotNil(t,m)
	err := m.Close(tid.NewString(0,"12"),"fixed")
	assert.NotNil(t,err)
	_, ok := err.(net.Error);
	assert.True(t,ok)
}
