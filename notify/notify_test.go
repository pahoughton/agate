/* 2019-02-17 (cc) <paul4hough@gmail.com>
   validate notify methods
*/
package notify

import (
	"strconv"
	"testing"
	"net/http/httptest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/pahoughton/agate/config"

	"github.com/pahoughton/agate/notify/mock"
	"github.com/pahoughton/agate/notify/gitlab"
	"github.com/pahoughton/agate/notify/hpsm"
)

type TestNSysSrv struct {
	mock *mock.MockServer
	hpsm *hpsm.MockServer
	gitlab *gitlab.MockServer
}

func TestGroup(t *testing.T) {
	cfg := config.New()
	obj := New(cfg.Notify,false)
	assert.NotNil(t,obj)
	got := obj.Group(NSysMock)
	assert.Equal(t,got,"")
	obj.Del()

	exp := "agate-test"
	cfg.Notify.Sys.Gitlab.Group = exp
	cfg.Notify.Sys.Hpsm.Group = exp + "hpsm"
	obj = New(cfg.Notify,false)
	assert.NotNil(t,obj)
	assert.Equal(t,exp,obj.Group(NSysGitlab))
	assert.Equal(t,exp + "hpsm",obj.Group(NSysHpsm))
	obj.Del()
}

func TestCrud(t *testing.T) {
	nsys := TestNSysSrv{
		hpsm: &hpsm.MockServer{},
		mock: &mock.MockServer{},
		gitlab: gitlab.NewMockServer(),
	}

	cfg := config.New()

	svcmock := httptest.NewServer(nsys.mock)
	defer svcmock.Close()
	cfg.Notify.Sys.Mock.Url = svcmock.URL

	svcgitlab := httptest.NewServer(nsys.gitlab)
	defer svcgitlab.Close()
	cfg.Notify.Sys.Gitlab.Url = svcgitlab.URL

	svchpsm := httptest.NewServer(nsys.hpsm)
	defer svchpsm.Close()
	cfg.Notify.Sys.Hpsm.Url = svchpsm.URL

	obj := New(cfg.Notify,false)
	defer obj.Del()
	assert.NotNil(t,obj)

	expHits := nsys.mock.Hits
	grp := "group"
	title := "note title"
	desc := "test node desc"
	nid, err := obj.Create(NSysMock,grp,title,desc,false,false)
	require.Nil(t,err)
	require.NotNil(t,nid)
	assert.Equal(t,strconv.Itoa(nsys.mock.Nid),nid.Id())
	expHits += 1
	assert.Equal(t,expHits,nsys.mock.Hits)
	assert.Nil(t,obj.Update(nid,desc))
	expHits += 1
	assert.Equal(t,expHits,nsys.mock.Hits)
	assert.Nil(t,obj.Close(nid,desc))
	expHits += 2
	assert.Equal(t,expHits,nsys.mock.Hits)


	expHits = nsys.gitlab.Hits
	nid, err = obj.Create(NSysGitlab,grp,title,desc,false,false)
	require.Nil(t,err)
	assert.NotNil(t,nid)
	expHits += 1
	assert.Equal(t,expHits,nsys.gitlab.Hits)
	assert.Nil(t,obj.Update(nid,desc))
	expHits += 1
	assert.Equal(t,expHits,nsys.gitlab.Hits)
	assert.Nil(t,obj.Close(nid,desc))
	expHits += 2
	assert.Equal(t,expHits,nsys.gitlab.Hits)

	expHits = nsys.hpsm.Hits
	nid, err = obj.Create(NSysHpsm,grp,title,desc,false,false)
	require.Nil(t,err)
	assert.NotNil(t,nid)
	expHits += 1
	assert.Equal(t,expHits,nsys.hpsm.Hits)
	assert.Nil(t,obj.Update(nid,desc))
	expHits += 1
	assert.Equal(t,expHits,nsys.hpsm.Hits)
	assert.Nil(t,obj.Close(nid,desc))
	expHits += 1
	assert.Equal(t,expHits,nsys.hpsm.Hits)
}
