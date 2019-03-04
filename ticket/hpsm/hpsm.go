/* 2019-01-19 (cc) <paul.houghton.ywi9@statefarm.com>
   HPSM ticket interface
*/

package hpsm

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pahoughton/agate/config"
	"github.com/pahoughton/agate/ticket/tid"
)

const (
	In2SoapActBase = `http://sf/application/automation/ws/sm/im-v2/so/`
	In3SoapActBase = `http://sf/application/automation/ws/sm/im-v3/so/`

	StatusSUCCESS = "SUCCESS"

	StatusFAILURE = "FAILURE"

	XmlNsSoap = "http://schemas.xmlsoap.org/soap/envelope/"

	XmlNsSO2 = "http://sf/application/automation/ws/sm/im-v2_0/so"
	XmlNsTO2 = "http://sf/application/automation/ws/sm/im-v2/to"

	XmlNsSO3 = "http://sf/application/automation/ws/sm/im-v3_0/so"
	XmlNsTO3 = "http://sf/application/automation/ws/sm/im-v3/to"

)

type Hpsm struct {
	tsys		uint8
	grp			string
	debug		bool
	BaseUrl		string
	CreateEp	string
	UpdateEp	string
	CloseEp		string
	User		string
	Pass		string
	Defaults	map[string]string
}

type ReqEnvelope struct {
	XMLName		xml.Name	`xml:"soapenv:Envelope"`
	XMLNsSoap	string		`xml:"xmlns:soapenv,attr"`
	XMLNsSO		string		`xml:"xmlns:so,attr"`
	XMLNsTO		string		`xml:"xmlns:to,attr"`

	Header  ReqHeader
	Body    ReqBody
}

type ReqHeader struct {
	XMLName xml.Name `xml:"soapenv:Header"`

	Items []interface{} `xml:",omitempty"`
}

type ReqBody struct {
	XMLName xml.Name `xml:"soapenv:Body"`

	Content interface{} `xml:",omitempty"`
}

type CreateIncidentRequest struct {
	XMLName		xml.Name	`xml:"so:createIncidentRequest"`
	Incident	In2ReqIncident
}

type UpdateIncidentRequest struct {
	XMLName		xml.Name	`xml:"so:updateIncidentRequest"`
	Incident	In2ReqIncident
}

type CloseIncidentRequest struct {
	XMLName		xml.Name	`xml:"so:CloseIncidentRequest"`
	Incident	Incident3
}

type In2RespEnvelope struct {
	XMLName		xml.Name	`xml:"Envelope"`

	Header  RespHeader
	Body    In2RespBody
}

type RespHeader struct {
	XMLName xml.Name `xml:"Header"`

	Items []interface{} `xml:",omitempty"`
}

type In2RespBody struct {
	XMLName xml.Name `xml:"Body"`

	Resp	In2Resp
}

type In2Resp struct {
	XMLName xml.Name `xml:"incidentResponse"`

	Incident		In2RespIncident
	StatusMessage	*StatusMessage
}


type In3RespEnvelope struct {
	XMLName		xml.Name	`xml:"Envelope"`

	Header  RespHeader
	Body    In3RespBody
}

type In3RespBody struct {
	XMLName xml.Name `xml:"Body"`

	Resp	In3Resp
}

type In3Resp struct {
	XMLName xml.Name `xml:"IncidentResponse"`

	Incident		In3RespIncident
	StatusMessage	*StatusMessage
}

type StatusMessage struct {

	XMLName xml.Name `xml:"StatusMessage"`

	Status *string `xml:"status,omitempty"`

	Message []string `xml:"message,omitempty"`
}

type In2ReqIncident struct {
	XMLName		xml.Name	`xml:"so:Incident"`

	AffectedCI			string `xml:"to:affectedCI,omitempty"`
	AssignmentGroup		string `xml:"to:assignmentGroup,omitempty"`
	BriefDescription	string `xml:"to:briefDescription,omitempty"`
	Category			string `xml:"to:category,omitempty"`
	CurrentUpdate		string `xml:"to:currentUpdate,omitempty"`
	Customer			string `xml:"to:customer,omitempty"`
	Impact				int32  `xml:"to:impact,omitempty"`
	IncidentDescription	string `xml:"to:incidentDescription,omitempty"`
	IncidentID			string `xml:"to:incidentID,omitempty"`
	Priority			int32  `xml:"to:priority,omitempty"`
	Service				string `xml:"to:service,omitempty"`
	Subcategory			string `xml:"to:subcategory,omitempty"`
	Type_				string `xml:"to:type,omitempty"`
	Urgency				int32  `xml:"to:urgency,omitempty"`
}

type In2RespIncident struct {
	XMLName		xml.Name	`xml:"Incident"`

	IncidentID	string		`xml:"incidentID,omitempty"`
}

type Incident3 struct {
	XMLName		xml.Name	`xml:"http://sf/application/automation/ws/sm/im-v3/to Incident"`

	Assignee	string `xml:"Assignee,omitempty"`
	ClosureCode	string `xml:"ClosureCode,omitempty"`
	IncidentID	string `xml:"IncidentID,omitempty"`
	Solution	string `xml:"Solution,omitempty"`
	Status		string `xml:"Status,omitempty"`
}

type In3RespIncident struct {
	XMLName		xml.Name	`xml:"Incident"`

	IncidentID	string		`xml:"IncidentID,omitempty"`
}

func New(cfg config.TSysHpsm, tsys int, dbg bool) *Hpsm {
	h := &Hpsm{
		tsys:		uint8(tsys),
		debug:		dbg,
		grp:		cfg.Group,
		BaseUrl:	cfg.Url,
		CreateEp:	cfg.CreateEp,
		UpdateEp:	cfg.UpdateEp,
		CloseEp:	cfg.CloseEp,
		User:		cfg.User,
		Pass:		cfg.Pass,
		Defaults:	cfg.Defaults,
	}
	return h
}

func (h *Hpsm) Group() string {
	return h.grp
}

func (h *Hpsm) GoodStatusMesg(sm *StatusMessage) error {
	if sm == nil ||
		sm.Status == nil ||
		*sm.Status != StatusSUCCESS {
		return  fmt.Errorf("resp status:  %v\n",sm)
	}
	return nil
}

func (h *Hpsm) PostSoap(url, sact string, reqObj, resp interface{}) error {

	reqXml, err := xml.MarshalIndent(reqObj,"  ","  ")
	if err != nil {
		return fmt.Errorf("hpsm marshal req %v",err)
	}

	if h.debug {
		fmt.Println("hpsm req xml: "+string(reqXml))
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(reqXml))
	if err != nil {
		return err
	}

	req.SetBasicAuth(h.User, h.Pass)

	req.Header.Add("Content-Type", "text/xml; charset=\"utf-8\"")
	req.Header.Add("SOAPAction", sact)
	req.Header.Set("User-Agent", "agate")

	req.Close = true

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client := &http.Client{Transport: tr}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	rawbody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if len(rawbody) == 0 {
		return fmt.Errorf("hpsm post no response body")
	}

	if h.debug {
		fmt.Println("hpsm resp: " + string(rawbody))
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("hpsm post resp status %d - %s\n",
			res.StatusCode, string(rawbody))
	}

	err = xml.Unmarshal(rawbody, resp)
	if err != nil {
		return fmt.Errorf("hpsm post resp unmarshal: %v",err)
	}

	if h.debug {
		fmt.Printf("hpsm resp unmarsh: %v\n",resp)
	}
	return nil

}

func (h *Hpsm) Create(wg, title, desc string) (tid.Tid, error) {

	ir := &CreateIncidentRequest{
		Incident:	In2ReqIncident{
			AffectedCI:				"Infrastructure",
			AssignmentGroup:		wg,
			BriefDescription:		title,
			IncidentDescription:	desc,
			Category:				"Incident",
			Customer:				"ip_soft_int",
			Impact:					4,
			Priority:				4,
			Service:				"EVENTING-EXCEPTION",
			Subcategory:			"Failure",
			Type_:					"Error Message",
			Urgency:				4,
		},
	}

	url := h.BaseUrl + "/" + h.CreateEp
	act := In2SoapActBase + "createIncident"

	reqEnv := ReqEnvelope{}
	reqEnv.XMLNsSoap = XmlNsSoap
	reqEnv.XMLNsSO = XmlNsSO2
	reqEnv.XMLNsTO = XmlNsTO2

	reqEnv.Body.Content = ir

	respEnv := new(In2RespEnvelope)

	if err := h.PostSoap(url,act,reqEnv,respEnv); err != nil {
		return nil, err
	}
	if err := h.GoodStatusMesg(respEnv.Body.Resp.StatusMessage); err != nil {
		return nil, err
	}
	if h.debug {
		fmt.Println("hpsm create ID: "+respEnv.Body.Resp.Incident.IncidentID)
	}
	return tid.NewString(h.tsys,respEnv.Body.Resp.Incident.IncidentID), nil

}

func (h *Hpsm)Update(id tid.Tid, cmt string) error {

	ir := UpdateIncidentRequest{
		Incident:	In2ReqIncident{
			CurrentUpdate:		cmt,
			IncidentID:			id.String(),
		},
	}

	url := h.BaseUrl + "/" + h.UpdateEp
	act := In2SoapActBase + "updateIncident"

	reqEnv := ReqEnvelope{}
	reqEnv.XMLNsSoap = XmlNsSoap
	reqEnv.XMLNsSO = XmlNsSO2
	reqEnv.XMLNsTO = XmlNsTO2

	reqEnv.Body.Content = ir

	respEnv := new(In2RespEnvelope)

	if err := h.PostSoap(url,act,reqEnv,respEnv); err != nil {
		return err
	}
	if err := h.GoodStatusMesg(respEnv.Body.Resp.StatusMessage); err != nil {
		return err
	}
	if h.debug {
		fmt.Println("hpsm update ID: "+respEnv.Body.Resp.Incident.IncidentID)
	}
	return nil
}

func (h *Hpsm)Close(id tid.Tid, cmt string) error {

	ir := CloseIncidentRequest{
		Incident:	Incident3{
			Assignee:		"ip_soft_int",
			ClosureCode:	"Automatically Closed",
			IncidentID:		id.String(),
			Solution:		cmt,
			Status:			"Resolved",
		},
	}

	url := h.BaseUrl + "/" + h.CloseEp
	act := In3SoapActBase + "CloseIncident"

	reqEnv := ReqEnvelope{}
	reqEnv.XMLNsSoap = XmlNsSoap
	reqEnv.XMLNsSO = XmlNsSO2
	reqEnv.XMLNsTO = XmlNsTO2

	reqEnv.Body.Content = ir

	respEnv := new(In3RespEnvelope)

	if err := h.PostSoap(url,act,reqEnv,respEnv); err != nil {
		return err
	}
	if err := h.GoodStatusMesg(respEnv.Body.Resp.StatusMessage); err != nil {
		return err
	}
	if h.debug {
		fmt.Println("hpsm close ID: "+respEnv.Body.Resp.Incident.IncidentID)
	}
	return nil
}
