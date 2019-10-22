/* 2019-02-17 (cc) <paul4hough@gmail.com>
   validate hpsm ticket interface
*/
package hpsm

import (
	"net"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pahoughton/agate/config"
)

func tcfg() config.NSysHpsm {
	return config.NSysHpsm{
		Url: "http://hpsm/api",
		User: "user",
		Pass: "secret-sauce",
		CreateEp: "incident2",
		UpdateEp: "incident2",
		CloseEp: "incident3",
	}
}
func TestNew(t *testing.T) {
	h := New(tcfg(),1,false)
	assert.NotNil(t,h)
}

func TestGroup(t *testing.T) {
	exp := "WG1234"
	h := New(config.NSysHpsm{Group: exp},1,false)
	assert.Equal(t,exp,h.Group())
}

func TestCreate(t *testing.T) {

	mock := &MockServer{Nid: 1233}
	ms := httptest.NewServer(mock)
	defer ms.Close()

	cfg := tcfg()
	cfg.Url = ms.URL
	h := New(cfg,1,false)

	nid, err := h.Create("storage","disk full","disk is full")

	assert.NotNil(t,nid)
	assert.Nil(t,err)
	assert.Equal(t,"IM1234",nid.Id())
}
func TestCreateNetError(t *testing.T) {
	cfg := tcfg()
	cfg.Url = "http://localhost:31001/hpsm"
	h := New(cfg,1,false)

	nid, err := h.Create("storage","disk full","disk is full")

	assert.Nil(t,nid)
	assert.NotNil(t,err)
	_, ok := err.(net.Error);
	assert.True(t,ok)
}

func TestUpdate(t *testing.T) {
	mock := &MockServer{Nid: 1233}
	ms := httptest.NewServer(mock)
	defer ms.Close()

	cfg := tcfg()
	cfg.Url = ms.URL
	h := New(cfg,1,false)

	err := h.Update(nid.NewString(h.tsys,"IM1234"),"disk still full")

	assert.Nil(t,err)
}

func TestUpdateNetError(t *testing.T) {
	cfg := tcfg()
	cfg.Url = "http://localhost:31001/hpsm"
	h := New(cfg,1,false)

	err := h.Update(nid.NewString(1,"IM1234"),"disk still full")

	assert.NotNil(t,err)
	_, ok := err.(net.Error);
	assert.True(t,ok)
}

func TestClose(t *testing.T) {
	mock := &MockServer{Nid: 1233}
	ms := httptest.NewServer(mock)
	defer ms.Close()

	cfg := tcfg()
	cfg.Url = ms.URL
	h := New(cfg,1,false)

	err := h.Close(nid.NewString(1,"IM1234"),"fixed")

	assert.Nil(t,err)
}

func TestCloseNetError(t *testing.T) {
	cfg := tcfg()
	cfg.Url = "http://localhost:31001/hpsm"
	h := New(cfg,1,false)

	err := h.Close(nid.NewString(1,"IM1234"),"fixed")

	assert.NotNil(t,err)
	_, ok := err.(net.Error);
	assert.True(t,ok)
}
