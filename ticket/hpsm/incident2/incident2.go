package incident2

import (
	"encoding/xml"
	"github.com/hooklift/gowsdl/soap"
	"time"
)

// against "unused imports"
var _ time.Time
var _ xml.Name

type CreateIncidentRequest struct {
	XMLName xml.Name `xml:"http://sf/application/automation/ws/sm/im-v2_0/so createIncidentRequest"`

	Incident *Incident `xml:"Incident,omitempty"`
}

type UpdateIncidentRequest struct {
	XMLName xml.Name `xml:"http://sf/application/automation/ws/sm/im-v2_0/so updateIncidentRequest"`

	Incident *Incident `xml:"Incident,omitempty"`
}

type RetrieveIncidentRequest struct {
	XMLName xml.Name `xml:"http://sf/application/automation/ws/sm/im-v2_0/so retrieveIncidentRequest"`

	Incident *Incident `xml:"Incident,omitempty"`
}

type CreateMultipleIncidentsRequest struct {
	XMLName xml.Name `xml:"http://sf/application/automation/ws/sm/im-v2_0/so createMultipleIncidentsRequest"`

	Incidents []*Incident `xml:"Incidents,omitempty"`
}

type UpdateMultipleIncidentsRequest struct {
	XMLName xml.Name `xml:"http://sf/application/automation/ws/sm/im-v2_0/so updateMultipleIncidentsRequest"`

	Incidents []*Incident `xml:"Incidents,omitempty"`
}

type RetrieveMultipleIncidentsRequest struct {
	XMLName xml.Name `xml:"http://sf/application/automation/ws/sm/im-v2_0/so retrieveMultipleIncidentsRequest"`

	Incidents []*Incident `xml:"Incidents,omitempty"`
}

type CountIncidentsRequest struct {
	XMLName xml.Name `xml:"http://sf/application/automation/ws/sm/im-v2_0/so countIncidentsRequest"`

	IncidentQuery *IncidentQuery `xml:"IncidentQuery,omitempty"`
}

type FindIncidentsRequest struct {
	XMLName xml.Name `xml:"http://sf/application/automation/ws/sm/im-v2_0/so findIncidentsRequest"`

	IncidentQuery *IncidentQuery `xml:"IncidentQuery,omitempty"`

	PageInfo *PagingData `xml:"PageInfo,omitempty"`
}

type CreateRecoveryTaskRequest struct {
	XMLName xml.Name `xml:"http://sf/application/automation/ws/sm/im-v2_0/so createRecoveryTaskRequest"`

	ParentIncidentID string `xml:"ParentIncidentID,omitempty"`

	RecoveryTask *RecoveryTask `xml:"RecoveryTask,omitempty"`
}

type UpdateRecoveryTaskRequest struct {
	XMLName xml.Name `xml:"http://sf/application/automation/ws/sm/im-v2_0/so updateRecoveryTaskRequest"`

	RecoveryTask *RecoveryTask `xml:"RecoveryTask,omitempty"`
}

type RetrieveRecoveryTaskRequest struct {
	XMLName xml.Name `xml:"http://sf/application/automation/ws/sm/im-v2_0/so retrieveRecoveryTaskRequest"`

	RecoveryTask *RecoveryTask `xml:"RecoveryTask,omitempty"`
}

type CreateMultipleRecoveryTasksRequest struct {
	XMLName xml.Name `xml:"http://sf/application/automation/ws/sm/im-v2_0/so createMultipleRecoveryTasksRequest"`

	ParentIncidentID string `xml:"ParentIncidentID,omitempty"`

	RecoveryTasks []*RecoveryTask `xml:"RecoveryTasks,omitempty"`
}

type UpdateMultipleRecoveryTasksRequest struct {
	XMLName xml.Name `xml:"http://sf/application/automation/ws/sm/im-v2_0/so updateMultipleRecoveryTasksRequest"`

	RecoveryTasks []*RecoveryTask `xml:"RecoveryTasks,omitempty"`
}

type RetrieveMultipleRecoveryTasksRequest struct {
	XMLName xml.Name `xml:"http://sf/application/automation/ws/sm/im-v2_0/so retrieveMultipleRecoveryTasksRequest"`

	RecoveryTasks []*RecoveryTask `xml:"RecoveryTasks,omitempty"`
}

type CountRecoveryTasksRequest struct {
	XMLName xml.Name `xml:"http://sf/application/automation/ws/sm/im-v2_0/so countRecoveryTasksRequest"`

	RecoveryTaskQuery *RecoveryTaskQuery `xml:"RecoveryTaskQuery,omitempty"`
}

type FindRecoveryTasksRequest struct {
	XMLName xml.Name `xml:"http://sf/application/automation/ws/sm/im-v2_0/so findRecoveryTasksRequest"`

	RecoveryTaskQuery *RecoveryTaskQuery `xml:"RecoveryTaskQuery,omitempty"`

	PageInfo *PagingData `xml:"PageInfo,omitempty"`
}

type FindRelatedTasksRequest struct {
	XMLName xml.Name `xml:"http://sf/application/automation/ws/sm/im-v2_0/so findRelatedTasksRequest"`

	IncidentID string `xml:"IncidentID,omitempty"`

	PageInfo *PagingData `xml:"PageInfo,omitempty"`
}

type IncidentResponse struct {
	XMLName xml.Name `xml:"http://sf/application/automation/ws/sm/im-v2_0/so incidentResponse"`

	Incident *Incident `xml:"Incident,omitempty"`

	StatusMessage *StatusMessage `xml:"StatusMessage,omitempty"`
}

type MultipleIncidentResponse struct {
	XMLName xml.Name `xml:"http://sf/application/automation/ws/sm/im-v2_0/so multipleIncidentResponse"`

	Responses []*IncidentResponseItem `xml:"Responses,omitempty"`

	StatusMessage *StatusMessage `xml:"StatusMessage,omitempty"`
}

type FindIncidentsResponse struct {
	XMLName xml.Name `xml:"http://sf/application/automation/ws/sm/im-v2_0/so findIncidentsResponse"`

	Incidents []*Incident `xml:"Incidents,omitempty"`

	PageInfo *PagingData `xml:"PageInfo,omitempty"`

	StatusMessage *StatusMessage `xml:"StatusMessage,omitempty"`
}

type RecoveryTaskResponse struct {
	XMLName xml.Name `xml:"http://sf/application/automation/ws/sm/im-v2_0/so recoveryTaskResponse"`

	RecoveryTask *RecoveryTask `xml:"RecoveryTask,omitempty"`

	StatusMessage *StatusMessage `xml:"StatusMessage,omitempty"`
}

type MultipleRecoveryTaskResponse struct {
	XMLName xml.Name `xml:"http://sf/application/automation/ws/sm/im-v2_0/so multipleRecoveryTaskResponse"`

	Responses []*RecoveryTaskResponseItem `xml:"Responses,omitempty"`

	StatusMessage *StatusMessage `xml:"StatusMessage,omitempty"`
}

type FindRecoveryTasksResponse struct {
	XMLName xml.Name `xml:"http://sf/application/automation/ws/sm/im-v2_0/so findRecoveryTasksResponse"`

	RecoveryTasks []*RecoveryTask `xml:"RecoveryTasks,omitempty"`

	PageInfo *PagingData `xml:"PageInfo,omitempty"`

	StatusMessage *StatusMessage `xml:"StatusMessage,omitempty"`
}

type CountResponse struct {
	XMLName xml.Name `xml:"http://sf/application/automation/ws/sm/im-v2_0/so countResponse"`

	Count int32 `xml:"Count,omitempty"`

	StatusMessage *StatusMessage `xml:"StatusMessage,omitempty"`
}

type Operation string

const (
	OperationEQUALS Operation = "EQUALS"

	OperationNOT_EQUAL_TO Operation = "NOT_EQUAL_TO"

	OperationSTARTS_WITH Operation = "STARTS_WITH"

	OperationDOES_NOT_START_WITH Operation = "DOES_NOT_START_WITH"

	OperationENDS_WITH Operation = "ENDS_WITH"

	OperationGREATER_THAN Operation = "GREATER_THAN"

	OperationGREATER_THAN_OR_EQUAL_TO Operation = "GREATER_THAN_OR_EQUAL_TO"

	OperationLESS_THAN Operation = "LESS_THAN"

	OperationLESS_THAN_OR_EQUAL_TO Operation = "LESS_THAN_OR_EQUAL_TO"

	OperationIS_NULL Operation = "IS_NULL"

	OperationIS_NOT_NULL Operation = "IS_NOT_NULL"

	OperationLIKE Operation = "LIKE"
)

type LogicalOperator string

const (
	LogicalOperatorAND LogicalOperator = "AND"

	LogicalOperatorOR LogicalOperator = "OR"
)

type ConditionType struct {
	Operator *LogicalOperator `xml:"Operator,omitempty"`

	Negate bool `xml:"Negate,omitempty"`
}

type PagingData struct {
	PageNumber int32 `xml:"PageNumber,omitempty"`

	PageSize int32 `xml:"PageSize,omitempty"`

	HasNextPage bool `xml:"HasNextPage,omitempty"`
}

type IntegerRange struct {
	Begin int32 `xml:"Begin,omitempty"`

	End int32 `xml:"End,omitempty"`
}

type NumberRange struct {
	Begin float64 `xml:"Begin,omitempty"`

	End float64 `xml:"End,omitempty"`
}

type DateRange struct {
	Begin time.Time `xml:"Begin,omitempty"`

	End time.Time `xml:"End,omitempty"`
}

type DateTimeRange struct {
	Begin time.Time `xml:"Begin,omitempty"`

	End time.Time `xml:"End,omitempty"`
}

type StringCriteria struct {
	Value []string `xml:"Value,omitempty"`

	Operation *Operation `xml:"Operation,omitempty"`
}

type StringArrayCriteria struct {
	Value []string `xml:"Value,omitempty"`

	Operation *Operation `xml:"Operation,omitempty"`
}

type IntegerCriteria struct {
	Value []int32 `xml:"Value,omitempty"`

	Operation *Operation `xml:"Operation,omitempty"`
}

type NumberCriteria struct {
	Value []float64 `xml:"Value,omitempty"`

	Operation *Operation `xml:"Operation,omitempty"`
}

type DateCriteria struct {
	Value time.Time `xml:"Value,omitempty"`

	Operation *Operation `xml:"Operation,omitempty"`
}

type DateTimeCriteria struct {
	Value time.Time `xml:"Value,omitempty"`

	Operation *Operation `xml:"Operation,omitempty"`
}

type BooleanCriteria struct {
	Value bool `xml:"Value,omitempty"`

	Operation *Operation `xml:"Operation,omitempty"`
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
	XMLName xml.Name `xml:"http://sf/application/automation/ws/response-v1/to wsNoRecordsFoundFaultMsg"`

	Messages []string `xml:"messages,omitempty"`

	CorrelationId string `xml:"correlationId,omitempty"`

	Actor *WsFaultActor `xml:"actor,omitempty"`
}

type StatusMessage struct {
	Status *Status `xml:"status,omitempty"`

	Message []string `xml:"message,omitempty"`

	Details *WsFault `xml:"details,omitempty"`
}

type Incident struct {
	AffectedCI string `xml:"affectedCI,omitempty"`

	Annotations string `xml:"annotations,omitempty"`

	Application string `xml:"application,omitempty"`

	Assignee string `xml:"assignee,omitempty"`

	AssignmentGroup string `xml:"assignmentGroup,omitempty"`

	BriefDescription string `xml:"briefDescription,omitempty"`

	Category string `xml:"category,omitempty"`

	CloseTime time.Time `xml:"closeTime,omitempty"`

	ClosedBy string `xml:"closedBy,omitempty"`

	CreatedFrom string `xml:"createdFrom,omitempty"`

	CurrentUpdate string `xml:"currentUpdate,omitempty"`

	Customer string `xml:"customer,omitempty"`

	ExpectedRecovery time.Time `xml:"expectedRecovery,omitempty"`

	Guid string `xml:"guid,omitempty"`

	Hostname string `xml:"hostname,omitempty"`

	Impact int32 `xml:"impact,omitempty"`

	ImpactEnd time.Time `xml:"impactEnd,omitempty"`

	ImpactStart time.Time `xml:"impactStart,omitempty"`

	IncidentDescription string `xml:"incidentDescription,omitempty"`

	IncidentID string `xml:"incidentID,omitempty"`

	IsMasterIncident bool `xml:"isMasterIncident,omitempty"`

	IsOutage bool `xml:"isOutage,omitempty"`

	KnowledgeItem string `xml:"knowledgeItem,omitempty"`

	Location string `xml:"location,omitempty"`

	ManagementServer string `xml:"managementServer,omitempty"`

	MasterIncidentID string `xml:"masterIncidentID,omitempty"`

	MessageGroup string `xml:"messageGroup,omitempty"`

	MessageText string `xml:"messageText,omitempty"`

	NextExpiration time.Time `xml:"nextExpiration,omitempty"`

	MonitoredObject string `xml:"monitoredObject,omitempty"`

	OpenTime time.Time `xml:"openTime,omitempty"`

	OpenedBy string `xml:"openedBy,omitempty"`

	Priority int32 `xml:"priority,omitempty"`

	ReasonForBreach string `xml:"reasonForBreach,omitempty"`

	ReasonForNoCI string `xml:"reasonForNoCI,omitempty"`

	RecommendedKI string `xml:"recommendedKI,omitempty"`

	RecordURL string `xml:"recordURL,omitempty"`

	RecoveredCI string `xml:"recoveredCI,omitempty"`

	RecoveryActions string `xml:"recoveryActions,omitempty"`

	Service string `xml:"service,omitempty"`

	ServiceOverride string `xml:"serviceOverride,omitempty"`

	SlaBreached bool `xml:"slaBreached,omitempty"`

	SourceID string `xml:"sourceID,omitempty"`

	Status string `xml:"status,omitempty"`

	Subcategory string `xml:"subcategory,omitempty"`

	Type_ string `xml:"type,omitempty"`

	UCmdbID string `xml:"uCmdbID,omitempty"`

	UpdateTime time.Time `xml:"updateTime,omitempty"`

	UpdateType string `xml:"updateType,omitempty"`

	UpdatedBy string `xml:"updatedBy,omitempty"`

	Urgency int32 `xml:"urgency,omitempty"`

	UserVariables string `xml:"userVariables,omitempty"`

	VendorInformation string `xml:"vendorInformation,omitempty"`

	VendorRecords []string `xml:"vendorRecords,omitempty"`

	Vendors []string `xml:"vendors,omitempty"`
}

type IncidentWhere struct {
	AffectedCICriteria *StringCriteria `xml:"affectedCICriteria,omitempty"`

	ApplicationCriteria *StringCriteria `xml:"applicationCriteria,omitempty"`

	AssigneeCriteria *StringCriteria `xml:"assigneeCriteria,omitempty"`

	AssignmentGroupCriteria *StringCriteria `xml:"assignmentGroupCriteria,omitempty"`

	BriefDescriptionCriteria *StringCriteria `xml:"briefDescriptionCriteria,omitempty"`

	CategoryCriteria *StringCriteria `xml:"categoryCriteria,omitempty"`

	CloseTimeCriteria *DateTimeCriteria `xml:"closeTimeCriteria,omitempty"`

	CloseTimeRange *DateTimeRange `xml:"closeTimeRange,omitempty"`

	ClosedByCriteria *StringCriteria `xml:"closedByCriteria,omitempty"`

	CreatedFromCriteria *StringCriteria `xml:"createdFromCriteria,omitempty"`

	CustomerCriteria *StringCriteria `xml:"customerCriteria,omitempty"`

	ExpectedRecoveryCriteria *DateTimeCriteria `xml:"expectedRecoveryCriteria,omitempty"`

	ExpectedRecoveryRange *DateTimeRange `xml:"expectedRecoveryRange,omitempty"`

	GuidCriteria *StringCriteria `xml:"guidCriteria,omitempty"`

	HostnameCriteria *StringCriteria `xml:"hostnameCriteria,omitempty"`

	ImpactCriteria *IntegerCriteria `xml:"impactCriteria,omitempty"`

	ImpactEndCriteria *DateTimeCriteria `xml:"impactEndCriteria,omitempty"`

	ImpactEndRange *DateTimeRange `xml:"impactEndRange,omitempty"`

	ImpactRange *IntegerRange `xml:"impactRange,omitempty"`

	ImpactStartCriteria *DateTimeCriteria `xml:"impactStartCriteria,omitempty"`

	ImpactStartRange *DateTimeRange `xml:"impactStartRange,omitempty"`

	IncidentIDCriteria *StringCriteria `xml:"incidentIDCriteria,omitempty"`

	IsMasterIncidentCriteria *BooleanCriteria `xml:"isMasterIncidentCriteria,omitempty"`

	IsOutageCriteria *BooleanCriteria `xml:"isOutageCriteria,omitempty"`

	KnowledgeItemCriteria *StringCriteria `xml:"knowledgeItemCriteria,omitempty"`

	LocationCriteria *StringCriteria `xml:"locationCriteria,omitempty"`

	ManagementServerCriteria *StringCriteria `xml:"managementServerCriteria,omitempty"`

	MasterIncidentIDCriteria *StringCriteria `xml:"masterIncidentIDCriteria,omitempty"`

	MessageGroupCriteria *StringCriteria `xml:"messageGroupCriteria,omitempty"`

	MessageTextCriteria *StringCriteria `xml:"messageTextCriteria,omitempty"`

	NextExpirationCriteria *DateTimeCriteria `xml:"nextExpirationCriteria,omitempty"`

	NextExpirationRange *DateTimeRange `xml:"nextExpirationRange,omitempty"`

	ObjectCriteria *StringCriteria `xml:"objectCriteria,omitempty"`

	OpenTimeCriteria *DateTimeCriteria `xml:"openTimeCriteria,omitempty"`

	OpenTimeRange *DateTimeRange `xml:"openTimeRange,omitempty"`

	OpenedByCriteria *StringCriteria `xml:"openedByCriteria,omitempty"`

	PriorityCriteria *IntegerCriteria `xml:"priorityCriteria,omitempty"`

	PriorityRange *IntegerRange `xml:"priorityRange,omitempty"`

	ReasonForBreachCriteria *StringCriteria `xml:"reasonForBreachCriteria,omitempty"`

	ReasonForNoCICriteria *StringCriteria `xml:"reasonForNoCICriteria,omitempty"`

	RecommendedKICriteria *StringCriteria `xml:"recommendedKICriteria,omitempty"`

	RecoveredCICriteria *StringCriteria `xml:"recoveredCICriteria,omitempty"`

	ServiceCriteria *StringCriteria `xml:"serviceCriteria,omitempty"`

	ServiceOverrideCriteria *StringCriteria `xml:"serviceOverrideCriteria,omitempty"`

	SlaBreachedCriteria *BooleanCriteria `xml:"slaBreachedCriteria,omitempty"`

	SourceIDCriteria *StringCriteria `xml:"sourceIDCriteria,omitempty"`

	StatusCriteria *StringCriteria `xml:"statusCriteria,omitempty"`

	SubcategoryCriteria *StringCriteria `xml:"subcategoryCriteria,omitempty"`

	TypeCriteria *StringCriteria `xml:"typeCriteria,omitempty"`

	UCmdbIDCriteria *StringCriteria `xml:"uCmdbIDCriteria,omitempty"`

	UpdateTimeCriteria *DateTimeCriteria `xml:"updateTimeCriteria,omitempty"`

	UpdateTimeRange *DateTimeRange `xml:"updateTimeRange,omitempty"`

	UpdatedByCriteria *StringCriteria `xml:"updatedByCriteria,omitempty"`

	UrgencyCriteria *IntegerCriteria `xml:"urgencyCriteria,omitempty"`

	UrgencyRange *IntegerRange `xml:"urgencyRange,omitempty"`
}

type IncidentQuery struct {
	QueryByExampleIncident *Incident `xml:"QueryByExampleIncident,omitempty"`

	IncidentCondition *IncidentCondition `xml:"IncidentCondition,omitempty"`
}

type IncidentCondition struct {
	Type *ConditionType `xml:"Type,omitempty"`

	Where []*IncidentWhere `xml:"Where,omitempty"`

	Condition []*IncidentCondition `xml:"Condition,omitempty"`
}

type IncidentResponseItem struct {
	Incident *Incident `xml:"Incident,omitempty"`

	StatusMessage *StatusMessage `xml:"StatusMessage,omitempty"`
}

type MultipleIncidentFault struct {
	XMLName xml.Name `xml:"http://sf/application/automation/ws/sm/im-v2/to wsIncidentMultipleTransactionFaultMsg"`

	Responses []*IncidentResponseItem `xml:"Responses,omitempty"`
}

type RecoveryTask struct {
	ActualEnd time.Time `xml:"actualEnd,omitempty"`

	ActualStart time.Time `xml:"actualStart,omitempty"`

	AffectedCI string `xml:"affectedCI,omitempty"`

	Assignee string `xml:"assignee,omitempty"`

	AssignmentGroup string `xml:"assignmentGroup,omitempty"`

	BriefDescription string `xml:"briefDescription,omitempty"`

	CloseTime time.Time `xml:"closeTime,omitempty"`

	ClosedBy string `xml:"closedBy,omitempty"`

	CurrentUpdate string `xml:"currentUpdate,omitempty"`

	Locations []string `xml:"locations,omitempty"`

	OpenTime time.Time `xml:"openTime,omitempty"`

	OpenedBy string `xml:"openedBy,omitempty"`

	ParentIncident string `xml:"parentIncident,omitempty"`

	Priority int32 `xml:"priority,omitempty"`

	RecordURL string `xml:"recordURL,omitempty"`

	ScheduledEnd time.Time `xml:"scheduledEnd,omitempty"`

	ScheduledStart time.Time `xml:"scheduledStart,omitempty"`

	Service string `xml:"service,omitempty"`

	Status string `xml:"status,omitempty"`

	TaskID string `xml:"taskID,omitempty"`

	TaskInstructions string `xml:"taskInstructions,omitempty"`

	UpdateType string `xml:"updateType,omitempty"`

	UpdatedBy string `xml:"updatedBy,omitempty"`
}

type RecoveryTaskWhere struct {
	ActualEndCriteria *DateTimeCriteria `xml:"actualEndCriteria,omitempty"`

	ActualEndRange *DateTimeRange `xml:"actualEndRange,omitempty"`

	ActualStartCriteria *DateTimeCriteria `xml:"actualStartCriteria,omitempty"`

	ActualStartRange *DateTimeRange `xml:"actualStartRange,omitempty"`

	AffectedCICriteria *StringCriteria `xml:"affectedCICriteria,omitempty"`

	AssigneeCriteria *StringCriteria `xml:"assigneeCriteria,omitempty"`

	AssignmentGroupCriteria *StringCriteria `xml:"assignmentGroupCriteria,omitempty"`

	BriefDescriptionCriteria *StringCriteria `xml:"briefDescriptionCriteria,omitempty"`

	CloseTimeCriteria *DateTimeCriteria `xml:"closeTimeCriteria,omitempty"`

	CloseTimeRange *DateTimeRange `xml:"closeTimeRange,omitempty"`

	ClosedByCriteria *StringCriteria `xml:"closedByCriteria,omitempty"`

	OpenTimeCriteria *DateTimeCriteria `xml:"openTimeCriteria,omitempty"`

	OpenTimeRange *DateTimeRange `xml:"openTimeRange,omitempty"`

	OpenedByCriteria *StringCriteria `xml:"openedByCriteria,omitempty"`

	ParentIncidentCriteria *StringCriteria `xml:"parentIncidentCriteria,omitempty"`

	PriorityCriteria *IntegerCriteria `xml:"priorityCriteria,omitempty"`

	PriorityRange *IntegerRange `xml:"priorityRange,omitempty"`

	ScheduledEndCriteria *DateTimeCriteria `xml:"scheduledEndCriteria,omitempty"`

	ScheduledEndRange *DateTimeRange `xml:"scheduledEndRange,omitempty"`

	ScheduledStartCriteria *DateTimeCriteria `xml:"scheduledStartCriteria,omitempty"`

	ScheduledStartRange *DateTimeRange `xml:"scheduledStartRange,omitempty"`

	ServiceCriteria *StringCriteria `xml:"serviceCriteria,omitempty"`

	StatusCriteria *StringCriteria `xml:"statusCriteria,omitempty"`

	TaskIDCriteria *StringCriteria `xml:"taskIDCriteria,omitempty"`

	UpdatedByCriteria *StringCriteria `xml:"updatedByCriteria,omitempty"`
}

type RecoveryTaskQuery struct {
	QueryByExampleRecoveryTask *RecoveryTask `xml:"QueryByExampleRecoveryTask,omitempty"`

	RecoveryTaskCondition *RecoveryTaskCondition `xml:"RecoveryTaskCondition,omitempty"`
}

type RecoveryTaskCondition struct {
	Type *ConditionType `xml:"Type,omitempty"`

	Where []*RecoveryTaskWhere `xml:"Where,omitempty"`

	Condition []*RecoveryTaskCondition `xml:"Condition,omitempty"`
}

type RecoveryTaskResponseItem struct {
	RecoveryTask *RecoveryTask `xml:"RecoveryTask,omitempty"`

	StatusMessage *StatusMessage `xml:"StatusMessage,omitempty"`
}

type MultipleRecoveryTaskFault struct {
	XMLName xml.Name `xml:"http://sf/application/automation/ws/sm/im-v2/to wsRecoveryTaskMultipleTransactionFaultMsg"`

	Responses []*RecoveryTaskResponseItem `xml:"Responses,omitempty"`
}

type IncidentManagement_v2_0 interface {

	// Error can be either of the following types:
	//
	//   - wsUnauthorizedUserException
	//   - wsGeneralUncheckedException
	//   - wsValidationFailedException

	CreateIncident(request *CreateIncidentRequest) (*IncidentResponse, error)

	// Error can be either of the following types:
	//
	//   - wsUnauthorizedUserException
	//   - wsGeneralUncheckedException
	//   - wsValidationFailedException
	//   - wsNoRecordsFoundException
	//   - wsResourceInUseException

	UpdateIncident(request *UpdateIncidentRequest) (*IncidentResponse, error)

	// Error can be either of the following types:
	//
	//   - wsUnauthorizedUserException
	//   - wsGeneralUncheckedException
	//   - wsNoRecordsFoundException

	RetrieveIncident(request *RetrieveIncidentRequest) (*IncidentResponse, error)

	// Error can be either of the following types:
	//
	//   - wsUnauthorizedUserException
	//   - wsGeneralUncheckedException
	//   - wsIncidentMultipleTransactionException

	CreateMultipleIncidents(request *CreateMultipleIncidentsRequest) (*MultipleIncidentResponse, error)

	// Error can be either of the following types:
	//
	//   - wsUnauthorizedUserException
	//   - wsGeneralUncheckedException
	//   - wsIncidentMultipleTransactionException

	UpdateMultipleIncidents(request *UpdateMultipleIncidentsRequest) (*MultipleIncidentResponse, error)

	// Error can be either of the following types:
	//
	//   - wsUnauthorizedUserException
	//   - wsGeneralUncheckedException
	//   - wsIncidentMultipleTransactionException

	RetrieveMultipleIncidents(request *RetrieveMultipleIncidentsRequest) (*MultipleIncidentResponse, error)

	// Error can be either of the following types:
	//
	//   - wsUnauthorizedUserException
	//   - wsGeneralUncheckedException
	//   - wsInvalidCriteriaException

	CountIncidents(request *CountIncidentsRequest) (*CountResponse, error)

	// Error can be either of the following types:
	//
	//   - wsUnauthorizedUserException
	//   - wsGeneralUncheckedException
	//   - wsInvalidCriteriaException
	//   - wsNoRecordsFoundException

	FindIncidents(request *FindIncidentsRequest) (*FindIncidentsResponse, error)

	// Error can be either of the following types:
	//
	//   - wsUnauthorizedUserException
	//   - wsGeneralUncheckedException
	//   - wsValidationFailedException

	CreateRecoveryTask(request *CreateRecoveryTaskRequest) (*RecoveryTaskResponse, error)

	// Error can be either of the following types:
	//
	//   - wsUnauthorizedUserException
	//   - wsGeneralUncheckedException
	//   - wsValidationFailedException
	//   - wsNoRecordsFoundException
	//   - wsResourceInUseException

	UpdateRecoveryTask(request *UpdateRecoveryTaskRequest) (*RecoveryTaskResponse, error)

	// Error can be either of the following types:
	//
	//   - wsUnauthorizedUserException
	//   - wsGeneralUncheckedException
	//   - wsNoRecordsFoundException

	RetrieveRecoveryTask(request *RetrieveRecoveryTaskRequest) (*RecoveryTaskResponse, error)

	// Error can be either of the following types:
	//
	//   - wsUnauthorizedUserException
	//   - wsGeneralUncheckedException
	//   - wsRecoveryTaskMultipleTransactionException

	CreateMultipleRecoveryTasks(request *CreateMultipleRecoveryTasksRequest) (*MultipleRecoveryTaskResponse, error)

	// Error can be either of the following types:
	//
	//   - wsUnauthorizedUserException
	//   - wsGeneralUncheckedException
	//   - wsRecoveryTaskMultipleTransactionException

	UpdateMultipleRecoveryTasks(request *UpdateMultipleRecoveryTasksRequest) (*MultipleRecoveryTaskResponse, error)

	// Error can be either of the following types:
	//
	//   - wsUnauthorizedUserException
	//   - wsGeneralUncheckedException
	//   - wsRecoveryTaskMultipleTransactionException

	RetrieveMultipleRecoveryTasks(request *RetrieveMultipleRecoveryTasksRequest) (*MultipleRecoveryTaskResponse, error)

	// Error can be either of the following types:
	//
	//   - wsUnauthorizedUserException
	//   - wsGeneralUncheckedException
	//   - wsInvalidCriteriaException

	CountRecoveryTasks(request *CountRecoveryTasksRequest) (*CountResponse, error)

	// Error can be either of the following types:
	//
	//   - wsUnauthorizedUserException
	//   - wsGeneralUncheckedException
	//   - wsInvalidCriteriaException
	//   - wsNoRecordsFoundException

	FindRecoveryTasks(request *FindRecoveryTasksRequest) (*FindRecoveryTasksResponse, error)

	// Error can be either of the following types:
	//
	//   - wsUnauthorizedUserException
	//   - wsGeneralUncheckedException
	//   - wsNoRecordsFoundException

	FindRelatedTasks(request *FindRelatedTasksRequest) (*FindRecoveryTasksResponse, error)
}

type incidentManagement_v2_0 struct {
	client *soap.Client
}

func NewIncidentManagement_v2_0(client *soap.Client) IncidentManagement_v2_0 {
	return &incidentManagement_v2_0{
		client: client,
	}
}

func (service *incidentManagement_v2_0) CreateIncident(request *CreateIncidentRequest) (*IncidentResponse, error) {
	response := new(IncidentResponse)
	err := service.client.Call("http://sf/application/automation/ws/sm/im-v2/so/createIncident", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *incidentManagement_v2_0) UpdateIncident(request *UpdateIncidentRequest) (*IncidentResponse, error) {
	response := new(IncidentResponse)
	err := service.client.Call("http://sf/application/automation/ws/sm/im-v2/so/updateIncident", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *incidentManagement_v2_0) RetrieveIncident(request *RetrieveIncidentRequest) (*IncidentResponse, error) {
	response := new(IncidentResponse)
	err := service.client.Call("http://sf/application/automation/ws/sm/im-v2/so/retrieveIncident", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *incidentManagement_v2_0) CreateMultipleIncidents(request *CreateMultipleIncidentsRequest) (*MultipleIncidentResponse, error) {
	response := new(MultipleIncidentResponse)
	err := service.client.Call("http://sf/application/automation/ws/sm/im-v2/so/createMultipleIncidents", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *incidentManagement_v2_0) UpdateMultipleIncidents(request *UpdateMultipleIncidentsRequest) (*MultipleIncidentResponse, error) {
	response := new(MultipleIncidentResponse)
	err := service.client.Call("http://sf/application/automation/ws/sm/im-v2/so/updateMultipleIncidents", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *incidentManagement_v2_0) RetrieveMultipleIncidents(request *RetrieveMultipleIncidentsRequest) (*MultipleIncidentResponse, error) {
	response := new(MultipleIncidentResponse)
	err := service.client.Call("http://sf/application/automation/ws/sm/im-v2/so/retrieveMultipleIncidents", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *incidentManagement_v2_0) CountIncidents(request *CountIncidentsRequest) (*CountResponse, error) {
	response := new(CountResponse)
	err := service.client.Call("http://sf/application/automation/ws/sm/im-v2/so/countIncidents", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *incidentManagement_v2_0) FindIncidents(request *FindIncidentsRequest) (*FindIncidentsResponse, error) {
	response := new(FindIncidentsResponse)
	err := service.client.Call("http://sf/application/automation/ws/sm/im-v2/so/findIncidents", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *incidentManagement_v2_0) CreateRecoveryTask(request *CreateRecoveryTaskRequest) (*RecoveryTaskResponse, error) {
	response := new(RecoveryTaskResponse)
	err := service.client.Call("http://sf/application/automation/ws/sm/im-v2/so/createRecoveryTask", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *incidentManagement_v2_0) UpdateRecoveryTask(request *UpdateRecoveryTaskRequest) (*RecoveryTaskResponse, error) {
	response := new(RecoveryTaskResponse)
	err := service.client.Call("http://sf/application/automation/ws/sm/im-v2/so/updateRecoveryTask", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *incidentManagement_v2_0) RetrieveRecoveryTask(request *RetrieveRecoveryTaskRequest) (*RecoveryTaskResponse, error) {
	response := new(RecoveryTaskResponse)
	err := service.client.Call("http://sf/application/automation/ws/sm/im-v2/so/retrieveRecoveryTask", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *incidentManagement_v2_0) CreateMultipleRecoveryTasks(request *CreateMultipleRecoveryTasksRequest) (*MultipleRecoveryTaskResponse, error) {
	response := new(MultipleRecoveryTaskResponse)
	err := service.client.Call("http://sf/application/automation/ws/sm/im-v2/so/createMultipleRecoveryTasks", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *incidentManagement_v2_0) UpdateMultipleRecoveryTasks(request *UpdateMultipleRecoveryTasksRequest) (*MultipleRecoveryTaskResponse, error) {
	response := new(MultipleRecoveryTaskResponse)
	err := service.client.Call("http://sf/application/automation/ws/sm/im-v2/so/updateMultipleRecoveryTasks", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *incidentManagement_v2_0) RetrieveMultipleRecoveryTasks(request *RetrieveMultipleRecoveryTasksRequest) (*MultipleRecoveryTaskResponse, error) {
	response := new(MultipleRecoveryTaskResponse)
	err := service.client.Call("http://sf/application/automation/ws/sm/im-v2/so/retrieveMultipleRecoveryTasks", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *incidentManagement_v2_0) CountRecoveryTasks(request *CountRecoveryTasksRequest) (*CountResponse, error) {
	response := new(CountResponse)
	err := service.client.Call("http://sf/application/automation/ws/sm/im-v2/so/countRecoveryTasks", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *incidentManagement_v2_0) FindRecoveryTasks(request *FindRecoveryTasksRequest) (*FindRecoveryTasksResponse, error) {
	response := new(FindRecoveryTasksResponse)
	err := service.client.Call("http://sf/application/automation/ws/sm/im-v2/so/findRecoveryTasks", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *incidentManagement_v2_0) FindRelatedTasks(request *FindRelatedTasksRequest) (*FindRecoveryTasksResponse, error) {
	response := new(FindRecoveryTasksResponse)
	err := service.client.Call("http://sf/application/automation/ws/sm/im-v2/so/findRelatedTasks", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
