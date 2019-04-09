/* 2019-02-23 (cc) <paul4hough@gmail.com>
   mock mock-ticket service for testing
*/
package hpsm

import (
	"fmt"
	"net/http"
	"strings"
)

type MockServer struct {
	Nid		int
	Hits	uint
}

func (s *MockServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	s.Hits += 1
	var body string
	if v, ok := r.Header["Soapaction"]; ! ok {
		w.WriteHeader(500)
		fmt.Fprintln(w,"missing soapaction header")
		return
	} else {
		sact := v[0]

		if strings.Index(sact,"create") > 0 {
			s.Nid += 1
			body = fmt.Sprintf(Resp2Xml,s.Nid)
		} else if strings.Index(sact,"Close") > 0 {
			body = fmt.Sprintf(Resp3Xml,s.Nid)
		} else {
			body = fmt.Sprintf(Resp2Xml,s.Nid)
		}
	}
	w.WriteHeader(200)
	if n, err := w.Write([]byte(body)); err == nil {
		if n != len(body) {
			panic("write len != data len")
		}
	}
}
const (
	Resp2Xml = `<?xml version="1.0" encoding="UTF-8"?>
<Envelope>
  <Header></Header>
  <Body>
    <incidentResponse>
      <Incident>
        <incidentID>IM%d</incidentID>
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
        <IncidentID>IM%d</IncidentID>
      </Incident>
      <StatusMessage>
        <status>SUCCESS</status>
      </StatusMessage>
    </IncidentResponse>
  </Body>
</Envelope>
`
)
