package toolsdk

import (
	"context"

	integrationsdk "github.com/domainry/domainry-integration-sdk"
	reportmodel "github.com/domainry/domainry-report-sdk/model"
)

// ConversationToolAuthorizer resolves current execution authority for a
// source-owned tool. Product hosts supply it; a tool implementation never
// infers authority from model input.
type ConversationToolAuthorizer func(context.Context, Request) (Authorization, error)

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
	MCP bool
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
