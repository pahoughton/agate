/* 2018-12-25 (cc) <paul4hough@gmail.com>
   mock ticket interface
*/
package mock

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pahoughton/agate/config"
	"github.com/pahoughton/agate/notify/note"
)

type Mock struct {
	name	string
	debug	bool
	url		string
}

func New(name string,cfg config.NSysMock,debug bool) *Mock {
	return &Mock{
		name:	name,
		debug:	debug,
		url:	cfg.Url,
	}
}

func (m *Mock)Group() string {
	return ""
}
func (self *Mock) Name() string {
	return self.name
}

func (m *Mock)Create(grp string, note note.Note, remcnt int) ([]byte, error) {

	tckt := map[string]string{
		"title":	note.Title(),
		"state":	"firing",
		"desc":		note.Desc(),
	}

	tcktJson, err := json.Marshal(tckt)
	if err != nil {
		panic(fmt.Errorf("json.Marshal: %s\n%+v\n",err.Error(),tckt))
	}

	resp, err := http.Post(
		m.url,
		"application/json",
		bytes.NewReader(tcktJson))

    if err != nil {
		return nil, err
    }
	defer resp.Body.Close()

	rcont, err := ioutil.ReadAll(resp.Body);
	if err != nil {
		return nil, err
	}

	var rmap map[string]string

	if err := json.Unmarshal(rcont, &rmap); err != nil {
		panic(err)
    }

	if resp.StatusCode != 200 {
		panic(fmt.Errorf("resp-status: %s\n%v",resp.Status,rcont))
	}

	id, ok := rmap["id"];
	if ! ok {
		panic(fmt.Errorf("no ticket id %v",rmap))
	}

	return []byte(id), nil
}

func (m *Mock)Update(note note.Note, cmt string) (bool,error) {

	tmap := map[string]string{
		"id": string(note.Nid),
		"comment": cmt,
	}
	tjson, err := json.Marshal(tmap)
	if err != nil {
		return false,fmt.Errorf("json.Marshal - %s",err.Error())
	}

	resp, err := http.Post(
		m.url,
		"application/json",
		bytes.NewReader(tjson))

	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return false, err
	}

	if resp.StatusCode != 200 {
		return false, fmt.Errorf("resp: "+resp.Status)
	}
	return false,nil
}

func (m *Mock)Close(note note.Note, cmt string) error {

	if len(cmt) > 0 {
		m.Update(note,cmt)
	}

	tmap := map[string]string{
		"id":		string(note.Nid),
		"state":	"closed",
	}
	tjson, err := json.Marshal(tmap)
	if err != nil {
		panic(fmt.Errorf("json.Marshal - %s",err.Error()))
	}

	resp, err := http.Post(
		m.url,
		"application/json",
		bytes.NewReader(tjson))

	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return err
	}

	rcont, err := ioutil.ReadAll(resp.Body);
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("resp: "+resp.Status+string(rcont))
	}
	return nil
}
