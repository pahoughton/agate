package incident3

import (
	"encoding/xml"
	"github.com/hooklift/gowsdl/soap"
	"time"
)

// against "unused imports"
var _ time.Time
var _ xml.Name

type CloseIncidentRequest struct {
	XMLName xml.Name `xml:"http://sf/application/automation/ws/sm/im-v3_0/so CloseIncidentRequest"`

	Incident *Incident `xml:"Incident,omitempty"`
}

type IncidentResponse struct {
	XMLName xml.Name `xml:"http://sf/application/automation/ws/sm/im-v3_0/so IncidentResponse"`

	Incident *Incident `xml:"Incident,omitempty"`

	StatusMessage *StatusMessage `xml:"StatusMessage,omitempty"`
}

type Incident struct {
	AffectedCI string `xml:"AffectedCI,omitempty"`

	Area string `xml:"Area,omitempty"`

	Assignee string `xml:"Assignee,omitempty"`

	AssignmentGroup string `xml:"AssignmentGroup,omitempty"`

	Category string `xml:"Category,omitempty"`

	CloseTime time.Time `xml:"CloseTime,omitempty"`

	ClosedBy string `xml:"ClosedBy,omitempty"`

	ClosureCode string `xml:"ClosureCode,omitempty"`

	Contact string `xml:"Contact,omitempty"`

	CreatedFrom string `xml:"CreatedFrom,omitempty"`

	CurrentUpdate string `xml:"CurrentUpdate,omitempty"`

	Description string `xml:"Description,omitempty"`

	Hostname string `xml:"Hostname,omitempty"`

	Impact int32 `xml:"Impact,omitempty"`

	IncidentID string `xml:"IncidentID,omitempty"`

	IsMasterIncident bool `xml:"IsMasterIncident,omitempty"`

	IsOutage bool `xml:"IsOutage,omitempty"`

	IsSlaBreached bool `xml:"IsSlaBreached,omitempty"`

	KnowledgeCandidate bool `xml:"KnowledgeCandidate,omitempty"`

	KnowledgeItem string `xml:"KnowledgeItem,omitempty"`

	Location string `xml:"Location,omitempty"`

	ManagementServer string `xml:"ManagementServer,omitempty"`

	MasterIncidentID string `xml:"MasterIncidentID,omitempty"`

	MessageText string `xml:"MessageText,omitempty"`

	NextExpiration time.Time `xml:"NextExpiration,omitempty"`

	OpenTime time.Time `xml:"OpenTime,omitempty"`

	OpenedBy string `xml:"OpenedBy,omitempty"`

	OutageEnd time.Time `xml:"OutageEnd,omitempty"`

	OutageStart time.Time `xml:"OutageStart,omitempty"`

	Priority int32 `xml:"Priority,omitempty"`

	ProblemCandidate bool `xml:"ProblemCandidate,omitempty"`

	ReasonForBreach string `xml:"ReasonForBreach,omitempty"`

	RecommendedKI string `xml:"RecommendedKI,omitempty"`

	Service string `xml:"Service,omitempty"`

	ServiceRecipient string `xml:"ServiceRecipient,omitempty"`

	Solution string `xml:"Solution,omitempty"`

	SourceID string `xml:"SourceID,omitempty"`

	Status string `xml:"Status,omitempty"`

	SubArea string `xml:"SubArea,omitempty"`

	Title string `xml:"Title,omitempty"`

	UpdateTime time.Time `xml:"UpdateTime,omitempty"`

	UpdateType string `xml:"UpdateType,omitempty"`

	UpdatedBy string `xml:"UpdatedBy,omitempty"`

	Urgency int32 `xml:"Urgency,omitempty"`

	Vendor string `xml:"Vendor,omitempty"`

	VendorTicket string `xml:"VendorTicket,omitempty"`
}

type Status string

const (
	StatusSUCCESS Status = "SUCCESS"

	StatusFAILURE Status = "FAILURE"
)

type WsFaultActor string

const (
	WsFaultActorSERVICE WsFaultActor = "SERVICE"

	WsFaultActorCLIENT WsFaultActor = "CLIENT"
)

type WsFault struct {
	XMLName xml.Name `xml:"http://sf/application/automation/ws/response-v1/to wsUnauthorizedUserFaultMsg"`

	Messages []string `xml:"messages,omitempty"`

	CorrelationId string `xml:"correlationId,omitempty"`

	Actor *WsFaultActor `xml:"actor,omitempty"`
}

type StatusMessage struct {
	Status *Status `xml:"status,omitempty"`

	Message []string `xml:"message,omitempty"`

	Details *WsFault `xml:"details,omitempty"`
}

type IncidentManagement_v3_0 interface {

	// Error can be either of the following types:
	//
	//   - WsUnauthorizedUserException
	//   - WsGeneralUncheckedException

	CloseIncident(request *CloseIncidentRequest) (*IncidentResponse, error)
}

type incidentManagement_v3_0 struct {
	client *soap.Client
}

func NewIncidentManagement_v3_0(client *soap.Client) IncidentManagement_v3_0 {
	return &incidentManagement_v3_0{
		client: client,
	}
}

func (service *incidentManagement_v3_0) CloseIncident(request *CloseIncidentRequest) (*IncidentResponse, error) {
	response := new(IncidentResponse)
	err := service.client.Call("http://sf/application/automation/ws/sm/im-v3/so/closeIncident", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
