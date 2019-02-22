/* 2019-02-17 (cc) <paul4hough@gmail.com>
   validate hpsm ticket interface
*/
package hpsm

import (
	"encoding/xml"
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

const (
	Resp2Xml = `<?xml version="1.0" encoding="UTF-8"?>
<Envelope>
  <Header></Header>
  <Body>
    <incidentResponse>
      <Incident>
        <incidentID>IM1234</incidentID>
      </Incident>
      <StatusMessage>
        <status>SUCCESS</status>
      </StatusMessage>
    </incidentResponse>
  </Body>
</Envelope>
`
	Resp3Xml = `<?xml version="1.0" encoding="UTF-8"?>
<Envelope>
  <Header></Header>
  <Body>
    <IncidentResponse>
      <Incident>
        <IncidentID>IM1234</IncidentID>
      </Incident>
      <StatusMessage>
        <status>SUCCESS</status>
      </StatusMessage>
    </IncidentResponse>
  </Body>
</Envelope>
`
)

func TestNew(t *testing.T) {
	h := New(
		config.TSysHpsm{
			Url: "http://hpsm/api",
			User: "user",
			Pass: "secret-sauce",
			CreateEp: "incident2",
			UpdateEp: "incident2",
			CloseEp: "incident3",
		},1,false)
	assert.NotNil(t,h)
}

func TestGroup(t *testing.T) {
	exp := "WG1234"
	h := New(
		config.TSysHpsm{
			Group: exp,
		},1,false)
	assert.Equal(t,exp,h.Group())
}

func TestCreate(t *testing.T) {

	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				cont, err := ioutil.ReadAll(r.Body);
				if err != nil {
					t.Error("readall")
					return
				}
				var dat []string
				if err := xml.Unmarshal(cont, &dat); err != nil {
					t.Errorf("xml.Unmarshal %v",err)
					return
				}
				fmt.Fprintln(w, Resp2Xml)
			}))
	defer ts.Close()

	h := New(
		config.TSysHpsm{
			Url: ts.URL,
			User: "user",
			Pass: "secret-sauce",
			CreateEp: "incident2",
			UpdateEp: "incident2",
			CloseEp: "incident3",
		},1,false)
	tid, err := h.Create("storage","disk full","disk is full")

	assert.NotNil(t,tid)
	assert.Nil(t,err)
	assert.Equal(t,"IM1234",tid.String())
}
func TestCreateNetError(t *testing.T) {
	h := New(
		config.TSysHpsm{
			Url: "http://localhost:31001/hpsm",
			User: "user",
			Pass: "secret-sauce",
			CreateEp: "incident2",
			UpdateEp: "incident2",
			CloseEp: "incident3",
		},1,false)
	tid, err := h.Create("storage","disk full","disk is full")
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
				var dat []string
				if err := xml.Unmarshal(cont, &dat); err != nil {
					t.Errorf("xml.Unmarshal %v",err)
					return
				}
				fmt.Fprintln(w, Resp2Xml)
			}))
	defer ts.Close()

	h := New(
		config.TSysHpsm{
			Url: ts.URL,
			User: "user",
			Pass: "secret-sauce",
			CreateEp: "incident2",
			UpdateEp: "incident2",
			CloseEp: "incident3",
		},1,false)
	err := h.Update(tid.NewString(h.tsys,"IM1234"),"disk still full")
	assert.Nil(t,err)
}

func TestUpdateNetError(t *testing.T) {
	h := New(
		config.TSysHpsm{
			Url: "http://localhost:31001/hpsm",
			User: "user",
			Pass: "secret-sauce",
			CreateEp: "incident2",
			UpdateEp: "incident2",
			CloseEp: "incident3",
		},1,false)
	err := h.Update(tid.NewString(1,"IM1234"),"disk still full")
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
				var dat []string
				if err := xml.Unmarshal(cont, &dat); err != nil {
					t.Errorf("xml.Unmarshal %v",err)
					return
				}
				fmt.Fprintln(w, Resp3Xml)
			}))
	defer ts.Close()

	h := New(
		config.TSysHpsm{
			Url: ts.URL,
			User: "user",
			Pass: "secret-sauce",
			CreateEp: "incident2",
			UpdateEp: "incident2",
			CloseEp: "incident3",
		},1,false)
	err := h.Close(tid.NewString(1,"IM1234"),"fixed")
	assert.Nil(t,err)
}

func TestCloseNetError(t *testing.T) {
	h := New(
		config.TSysHpsm{
			Url: "http://localhost:31001/hpsm",
			User: "user",
			Pass: "secret-sauce",
			CreateEp: "incident2",
			UpdateEp: "incident2",
			CloseEp: "incident3",
		},1,false)
	err := h.Close(tid.NewString(1,"IM1234"),"fixed")
	assert.NotNil(t,err)
	_, ok := err.(net.Error);
	assert.True(t,ok)
}
