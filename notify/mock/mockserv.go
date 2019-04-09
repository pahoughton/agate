/* 2019-02-23 (cc) <paul4hough@gmail.com>
   mock mock-ticket service for testing
*/
package mock

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

type MockServer struct {
	Nid		int
	Hits	uint
}

func (t *MockServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.Hits += 1
	cont, err := ioutil.ReadAll(r.Body)
	if err == nil {
		var dat map[string]string
		if err = json.Unmarshal(cont, &dat); err == nil {
			if _, ok := dat["id"]; ok == false {
				// new ticket
				t.Nid += 1
				resp := map[string]string{ "id": strconv.Itoa(t.Nid) }
				w.WriteHeader(200)
				if rjson, err := json.Marshal(resp); err == nil {
					if n, err := w.Write(rjson); err == nil {
						if n != len(rjson) {
							panic("write len != data len")
						}
						return
					} else {
						panic(err)
					}
				} else {
					panic(err)
				}
			} else {
				// update
				w.WriteHeader(200)
				return
			}
		}
	}
	w.WriteHeader(500)
	panic(err)

}
