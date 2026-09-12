// Package toolsdk defines tool contracts without depending on an Agent implementation.
package toolsdk

import (
	"context"
	"encoding/json"
	"time"
)

type Authority struct {
	Known       bool   `json:"known"`
	RuntimeID   string `json:"runtime_id"`
	WorkspaceID string `json:"workspace_id"`
	UserID      string `json:"user_id"`
	RoleKey     string `json:"role_key,omitempty"`
}

type Definition struct {
	Key            string          `json:"key"`
	Version        string          `json:"version"`
	Description    string          `json:"description"`
	InputSchema    json.RawMessage `json:"input_schema"`
	OutputSchema   json.RawMessage `json:"output_schema"`
	ActionKey      string          `json:"action_key"`
	Effect         string          `json:"effect"`      // read or write
	Idempotency    string          `json:"idempotency"` // natural, key or reconcile
	TimeoutMillis  int             `json:"timeout_ms"`
	MaxOutputBytes int             `json:"max_output_bytes"`
}

type Call struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Arguments string `json:"arguments"` // Complete JSON text; the executor validates the schema.
}

type Authorization struct {
	UserTimezone         string         `json:"user_timezone,omitempty"`
	Granted              bool           `json:"granted"`
	ConfirmationRequired bool           `json:"confirmation_required"`
	Revision             string         `json:"revision"`
	Evidence             map[string]any `json:"evidence,omitempty"`
}

type Confirmation struct {
	ID            string
	UserID        string
	ActionKey     string
	ToolVersion   string
	ArgumentsHash string
	ApprovedAt    time.Time
}

// ConfirmationVerifier is supplied by the execution owner. It verifies the
// persisted approval for this exact actor, run, call, definition and arguments.
// A structurally valid Confirmation value alone is never execution authority.
// Tools consume this port without importing the owner's implementation/store.
type ConfirmationVerifier interface {
	VerifyConversationToolConfirmation(context.Context, Request) (bool, error)
}

type Request struct {
	Authority      Authority
	ConversationID string
	RunID          string
	Step           int
	Call           Call
	Definition     Definition
	IdempotencyKey string
	ConfirmationID string
	Confirmation   *Confirmation
	// Server-only worker guard. Local transactional effects validate this
	// against persisted state; it is never accepted from model arguments.
	LeaseOwner string
	Fence      int64
}

type Result struct {
	// Source-owned disposition; never supplied by model arguments. An accepted
	// invocation may finish this tool call without finishing the business work.
	Completion string          `json:"completion,omitempty"`
	Status     string          `json:"status"` // completed, failed, pending or uncertain
	Content    json.RawMessage `json:"content,omitempty"`
	ErrorCode  string          `json:"error_code,omitempty"`
	ResourceID string          `json:"resource_id,omitempty"`
}

type Host interface {
	ConversationTools(context.Context, Authority) ([]Definition, error)
	AuthorizeConversationTool(context.Context, Request) (Authorization, error)
	InvokeConversationTool(context.Context, Request) (Result, error)
	ReconcileConversationTool(context.Context, Request) (Result, error)
}
type ResultAuthorizer interface {
	AuthorizeConversationToolResult(context.Context, Request, Result) error
}
