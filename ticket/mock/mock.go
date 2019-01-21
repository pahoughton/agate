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

	"github.com/pahoughton/agate/model"
)

type Mock struct {
	debug	bool
	url		string
}

func New(url string, debug bool) *Mock {
	m := &Mock{
		debug:	debug,
		url:	url,
	}
	return m
}

func (m *Mock)Create(a model.Alert) (string, error) {

	tckt := map[string]string{
		"title":	a.Title(),
		"state":	"firing",
		"desc":		a.Desc(),
	}

	tcktJson, err := json.Marshal(tckt)
	if err != nil {
		return "", fmt.Errorf("json.Marshal: %s\n%+v\n",err.Error(),tckt)
	}

	resp, err := http.Post(
		m.url,
		"application/json",
		bytes.NewReader(tcktJson))

    if err != nil {
		return "", fmt.Errorf("http.post-%s: %s \n%+v\n",
			m.url,err.Error(),tcktJson)
    }
	defer resp.Body.Close()

	rcont, err := ioutil.ReadAll(resp.Body);
	if err != nil {
		return "", err
	}

	var rmap map[string]string

	if err := json.Unmarshal(rcont, &rmap); err != nil {
		return "", err
    }

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("resp-status: %s\n%v",resp.Status,rcont)
	}

	tid, ok := rmap["id"];
	if ! ok {
		return "", fmt.Errorf("no ticket id %v",rmap)
	}

	return tid, nil
}

func (m *Mock)AddComment(tid string, cmt string) error {

	tmap := map[string]string{
		"id": tid,
		"comment": cmt,
	}
	tjson, err := json.Marshal(tmap)
	if err != nil {
		return fmt.Errorf("json.Marshal - %s",err.Error())
	}

	resp, err := http.Post(
		m.url,
		"application/json",
		bytes.NewReader(tjson))

	if err != nil {
		return fmt.Errorf("http.Post - %s",err.Error())
	}

	defer resp.Body.Close()

	rcont, err := ioutil.ReadAll(resp.Body);
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("resp: "+resp.Status+string(rcont))
	}
	return nil
}

func (m *Mock)Close(tid string) error {
	tmap := map[string]string{
		"id":		tid,
		"state":	"closed",
	}
	tjson, err := json.Marshal(tmap)
	if err != nil {
		return fmt.Errorf("json.Marshal - %s",err.Error())
	}

	resp, err := http.Post(
		m.url,
		"application/json",
		bytes.NewReader(tjson))

	if err != nil {
		return fmt.Errorf("http.Post - %s",err.Error())
	}

	defer resp.Body.Close()

	rcont, err := ioutil.ReadAll(resp.Body);
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("resp: "+resp.Status+string(rcont))
	}
	return nil
}
