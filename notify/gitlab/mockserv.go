/* 2019-02-23 (cc) <paul4hough@gmail.com>
   mock mock-ticket service for testing
*/
package gitlab

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

type MockServer struct {
	Nid		int
	Hits	uint
	Groups	map[int]string
}

func NewMockServer() *MockServer {
	return &MockServer{ Groups: make(map[int]string,16) }
}

func (t *MockServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	rtmpl := `{"id":1, "iid":%s, "title" : "Title of issue",
"description": "This is description of an issue",
"author" : {"id" : 1, "name": "snehal"}, "assignees":[{"id":1}]}`

	var resp string
	t.Hits += 1
	if path.Base(r.URL.String()) == "issues" {
		grp, _ := url.QueryUnescape(path.Base(path.Dir(r.URL.String())))
		t.Groups[t.Nid] = grp
		t.Nid += 1
		resp = fmt.Sprintf(rtmpl,strconv.Itoa(t.Nid))
	} else {
		resp = `{"id":1}`
	}

	w.WriteHeader(200)
	if n, err := w.Write([]byte(resp)); err == nil {
		if n != len(resp) {
			panic("write len != data len")
		} else {
			return
		}
	} else {
		panic(err)
	}
}
