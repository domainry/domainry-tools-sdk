package toolsdk

import "context"

// Availability is a live policy for registered, authorized tools. It composes
// user preferences and required connections without granting action permission.
// Return false for disabled, missing, expired or unknown required connections.
// Implementations honor ctx and perform no refresh or tool side effects.
type Availability interface {
	ConversationToolAvailable(context.Context, Authority, string) (bool, error)
}

// ResultReadAvailability keeps user preferences and source readiness while
// avoiding catalogs gated by permission to execute a tool. It grants no access
// to result contents; ResultReadAuthorizer remains mandatory for that path.
type ResultReadAvailability interface {
	ConversationToolResultReadAvailable(context.Context, Authority, string) (bool, error)
}

type ResultReadConnectionAvailability interface {
	ToolResultReadConnectionAvailable(context.Context, Authority, string) (bool, error)
}

// Catalog exposes only the current authorized deployment selection. A settings
// service must receive this before applying preferences, avoiding feedback from
// its own availability filter. Engines retain their Agent/Skill whitelists.
type Catalog interface {
	ConversationTools(context.Context, Authority) ([]Definition, error)
}

// ConnectionAvailability is implemented by trusted deployment composition.
// Account owners keep credentials and scope rules; Tools sees only readiness.
// A tool that needs no connection returns ready=true. An unknown mapping fails
// closed. This read must never refresh credentials or invoke a business action.
type ConnectionAvailability interface {
	ToolConnectionAvailable(context.Context, Authority, string) (bool, error)
}
