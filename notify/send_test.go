/* 2019-10-22 (cc) <paul4hough@gmail.com>
   FIXME what is this for?
*/
package notify

import (
	"testing"
	"net/http/httptest"
	"time"

	"github.com/stretchr/testify/assert"
	pmod "github.com/prometheus/common/model"

	"github.com/pahoughton/agate/config"
	"github.com/pahoughton/agate/notify/note"
	"github.com/pahoughton/agate/notify/mock"
)

func testNote() note.Note {
	return 	note.Note{
		From: "badthing",
		Labels: pmod.LabelSet{
			"alertname": "badthing",
			"instance": "box.nowhere.local:9001",
		},
		Alerts: []note.Alert{
			note.Alert{
				Name: "badthing",
				Labels: pmod.LabelSet{
					"alertname": "badthing",
					"instance": "box.nowhere.local:9001",
				},
				Starts: time.Now(),
			},
		},
	}
}


func TestSend(t *testing.T) {

	tsvc := &mock.MockServer{}
	svr := httptest.NewServer(tsvc)
	defer svr.Close()

	cfg := config.New()
	cfg.Notify.Sys.Mock.Url = svr.URL
	obj := cleanNew(cfg.Notify,"testdata",true)

	key  := Key{
		Sys: "mock",
		Grp: "alert",
		Key: []byte("key-01"),
	}
	note := testNote()

	hits := tsvc.Hits
	obj.Send(key,note,0)
	assert.True(t, tsvc.Hits > hits )

	// todo expand to include fail and retry validation
	obj.Del()
}
