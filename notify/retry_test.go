/* 2019-10-22 (cc) <paul4hough@gmail.com>
   FIXME what is this for?
*/
package notify

import (
	"testing"
	"net/http/httptest"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/pahoughton/agate/config"
	"github.com/pahoughton/agate/notify/mock"
)

func TestRetry(t *testing.T) {
	tsvc := &mock.MockServer{}
	svr := httptest.NewServer(tsvc)
	defer svr.Close()

	cfg := config.New()
	cfg.Notify.Sys.Mock.Url = svr.URL
	cfg.Notify.Retry = 0

	obj := cleanNew(cfg.Notify,"testdata",true)

	key  := Key{
		Sys: "mock",
		Grp: "alert",
		Key: []byte("key-01"),
	}
	note := testNote()


	hits := tsvc.Hits
	// w/o retry data
	obj.retryOnce()
	time.Sleep(20)
	assert.True(t, tsvc.Hits == hits )

	// now w/ a note in retry
	obj.retryMap.Store(key.KString(),retry{key,note,0})
	obj.retryOnce()
	time.Sleep(2)
	assert.True(t, tsvc.Hits > hits )
	_, found := obj.retryMap.Load(key.KString())
	assert.False(t, found)

	obj.Del()
}
