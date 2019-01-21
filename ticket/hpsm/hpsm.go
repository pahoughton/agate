/* 2019-01-19 (cc) <paul.houghton.ywi9@statefarm.com>
   HPSM ticket interface
*/

package hpsm

import (
	"fmt"

	"github.com/pahoughton/agate/model"

	in2 "github.com/pahoughton/agate/ticket/hpsm/incident2"
	in3 "github.com/pahoughton/agate/ticket/hpsm/incident3"

	"github.com/hooklift/gowsdl/soap"
)

type HPSM struct {
	Debug		bool
	In2			in2.IncidentManagement_v2_0
	In3			in3.IncidentManagement_v3_0
}

func New(url, user, pass string, dbg bool) *HPSM {
	h := &HPSM{
		Debug:		dbg,
		In2:		in2.NewIncidentManagement_v2_0(
			soap.NewClient( url,soap.WithBasicAuth(user,pass))),
		In3:		in3.NewIncidentManagement_v3_0(
			soap.NewClient( url,soap.WithBasicAuth(user,pass))),
	}
	return h
}

func (h *HPSM) Create(wg string, a model.Alert) (string, error) {

	ci := &in2.CreateIncidentRequest{
		Incident:	&in2.Incident{
			AssignmentGroup:		wg,
			BriefDescription:		a.Title(),
			IncidentDescription:	a.Desc(),
		},
	}

	ir, err := h.In2.CreateIncident(ci)
	if err != nil {
		return "", err
	}
	if ir.StatusMessage == nil ||
		ir.StatusMessage.Status == nil ||
		*ir.StatusMessage.Status != in2.StatusSUCCESS {
		return "", fmt.Errorf("create failed - %v\n",ir)
	}
	return ir.Incident.IncidentID, nil
}

func (h *HPSM)AddComment(tid string, cmt string) error {

	return fmt.Errorf("unsupported")
}

func (h *HPSM)Close(tid string) error {

	return fmt.Errorf("unsupported")
}
