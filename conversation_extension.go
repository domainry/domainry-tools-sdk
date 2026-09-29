package toolsdk

import (
	"context"
	"encoding/json"

	integrationsdk "github.com/domainry/domainry-integration-sdk"
	reportmodel "github.com/domainry/domainry-report-sdk/model"
)

// ConversationToolAuthorizer resolves current execution authority for a
// source-owned tool. Product hosts supply it; a tool implementation never
// infers authority from model input.
type ConversationToolAuthorizer func(context.Context, Request) (Authorization, error)

// ConversationBusinessRead is a closed source-owned business read request.
// Input is the canonical JSON payload for operation; implementations must
// reject unsupported operations and validate the complete payload before use.
type ConversationBusinessRead struct {
	Operation string          `json:"operation"`
	Input     json.RawMessage `json:"input"`
}

// ConversationBusinessEvidence binds one business read to the configured
// source, authenticated caller scope and exact returned data. It is persisted
// with a product tool result and never grants access by itself.
type ConversationBusinessEvidence struct {
	Version     int             `json:"version"`
	Source      string          `json:"source"`
	ScopeSHA256 string          `json:"scope_sha256"`
	Operation   string          `json:"operation"`
	Input       json.RawMessage `json:"input"`
	Data        json.RawMessage `json:"data"`
	HostProof   string          `json:"host_proof,omitempty"`
}

// ConversationBusinessSource is the narrow source-owner port used by
// product-specific business read tools. Runtime resolves current row/field
// access for every read and revalidates saved evidence before it is shown
// again; a Tools implementation never selects a principal or database.
type ConversationBusinessSource interface {
	BusinessSourceIdentity() string
	ReadConversationBusiness(context.Context, ConversationBusinessRead, Authority) (ConversationBusinessEvidence, error)
	RevalidateConversationBusiness(context.Context, ConversationBusinessEvidence, Authority) error
}

// ConversationBusinessAction is a closed product-owned mapping to one
// source-owned Runtime Action. The product tool supplies the business target
// and payload, while Runtime resolves and freezes the currently published
// Action contract; models never supply an Action execution version.
type ConversationBusinessAction struct {
	ObjectKey string          `json:"object_key"`
	ActionKey string          `json:"action_key"`
	RecordID  string          `json:"record_id,omitempty"`
	Data      json.RawMessage `json:"data"`
}

type ConversationBusinessRecordReference struct {
	ObjectKey string `json:"object_key"`
	RecordID  string `json:"record_id"`
}

// ConversationBusinessActionRequest contains only server-owned execution
// metadata in addition to the product's closed Action mapping. Arguments are
// the exact specialized-tool JSON covered by the persisted confirmation.
type ConversationBusinessActionRequest struct {
	Authority      Authority
	Action         ConversationBusinessAction
	ToolActionKey  string
	ToolVersion    string
	ConversationID string
	RunID          string
	CorrelationID  string
	Step           int
	CallID         string
	IdempotencyKey string
	Confirmation   *Confirmation
	Arguments      string
}

// ConversationBusinessActionResult is a source-owned acknowledgement. It
// intentionally contains effect references rather than unrestricted records.
// Completed results carry the exact Action evidence needed for current and
// historical result revalidation.
type ConversationBusinessActionResult struct {
	Status              string                                `json:"status"`
	ErrorCode           string                                `json:"error_code,omitempty"`
	InvocationID        string                                `json:"invocation_id,omitempty"`
	ObjectKey           string                                `json:"object_key"`
	ActionKey           string                                `json:"action_key"`
	RecordID            string                                `json:"record_id,omitempty"`
	CreatedRecords      []ConversationBusinessRecordReference `json:"created_records,omitempty"`
	UpdatedRecords      []ConversationBusinessRecordReference `json:"updated_records,omitempty"`
	DeletedRecords      []ConversationBusinessRecordReference `json:"deleted_records,omitempty"`
	RestoredRecords     []ConversationBusinessRecordReference `json:"restored_records,omitempty"`
	ReferenceCount      int                                   `json:"reference_count"`
	ReferencesTruncated bool                                  `json:"references_truncated,omitempty"`
	Evidence            *ConversationBusinessEvidence         `json:"evidence,omitempty"`
}

// ConversationBusinessActionSource is the optional mutation companion to
// ConversationBusinessSource. Authorization resolves the current Action
// contract and record scope. Reconciliation uses the same source receipt and
// must not treat an unknown outcome as permission to issue a second effect.
type ConversationBusinessActionSource interface {
	AuthorizeConversationBusinessAction(context.Context, ConversationBusinessAction, Authority) (Authorization, error)
	InvokeConversationBusinessAction(context.Context, ConversationBusinessActionRequest) (ConversationBusinessActionResult, error)
	ReconcileConversationBusinessAction(context.Context, ConversationBusinessActionRequest) (ConversationBusinessActionResult, error)
	RevalidateConversationBusinessAction(context.Context, ConversationBusinessEvidence, Authority) error
}

// ReportSource is the narrow Report SDK port consumed by the report_query
// tool implementation.
type ReportSource interface {
	BusinessSourceIdentity() string
	ReportCatalog(context.Context, reportmodel.ReportCatalogRequest, Authority) (reportmodel.ReportCatalog, error)
	QueryReport(context.Context, reportmodel.ReportObjectSQLRequest, Authority) (reportmodel.ReportQueryResult, error)
	AuthorizeReportResult(context.Context, reportmodel.ReportQueryResultAuthorization, Authority) error
}

// AnalysisSource is the narrow Report SDK port consumed by the analysis_run
// tool implementation.
type AnalysisSource interface {
	BusinessSourceIdentity() string
	AnalysisCatalog(context.Context, reportmodel.AnalysisCatalogRequest, Authority) (reportmodel.AnalysisCatalog, error)
	RunAnalysis(context.Context, reportmodel.AnalysisRequest, Authority) (reportmodel.AnalysisResult, error)
	AuthorizeAnalysisResult(context.Context, reportmodel.AnalysisResultAuthorization, Authority) error
}

// ConnectionAccounts is the public Integration SDK account catalog needed by
// account-backed conversation tools.
type ConnectionAccounts interface {
	ListConnectionAccounts(context.Context, integrationsdk.ConnectionAccountSubject) ([]integrationsdk.ConnectionAccount, error)
}

// ConnectionAccountSubjectResolver resolves the current caller and requested
// Integration action. Model arguments never enter this function.
type ConnectionAccountSubjectResolver func(context.Context, Authority, string) (integrationsdk.ConnectionAccountSubject, error)

// ConversationToolCapabilities contains topology facts only. It never carries
// stores or module implementations.
type ConversationToolCapabilities struct {
	Business bool
	MCP      bool
}

// ConversationMCPPorts are SDK-owned Integration ports supplied by the
// product host when MCP tools are enabled.
type ConversationMCPPorts struct {
	Accounts ConnectionAccounts
	Reads    integrationsdk.ConnectionAccountReads
	Writes   integrationsdk.ConnectionAccountWrites
	Subject  ConnectionAccountSubjectResolver
}

// ConversationToolAssembly contains only public owner ports. The selected
// Tools implementation assembles its private adapters around these ports.
type ConversationToolAssembly struct {
	Base           Host
	BusinessSource func() ConversationBusinessSource
	ReportSource   func() ReportSource
	AnalysisSource func() AnalysisSource
	Authorize      ConversationToolAuthorizer
	MCP            *ConversationMCPPorts
}

// ConversationToolFactory is selected by the outer product composition root.
// Runtime consumes this contract and never imports the Tools implementation.
type ConversationToolFactory interface {
	ConversationToolDefinitions(ConversationToolCapabilities) []Definition
	AssembleConversationTools(ConversationToolAssembly) (Host, error)
}
