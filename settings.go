package toolsdk

import "context"

// Preference is owned by one runtime/workspace/user and stable tool key. Role
// changes never transfer it. Revision zero means the default (enabled) setting.
type Preference struct {
	Key       string `json:"key"`
	Enabled   bool   `json:"enabled"`
	Revision  int64  `json:"revision"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

type ToolSetting struct {
	Preference
	Version     string `json:"version"`
	Description string `json:"description"`
	Available   bool   `json:"available"`
	State       string `json:"state"` // available, disabled, connection_unavailable, connection_unknown
}

// The current definition version and preference revision must both match.
// Browser values cannot nominate an owner, change a tool or grant permission.
type ToolSettingInput struct {
	Enabled          bool   `json:"enabled"`
	ExpectedRevision int64  `json:"expected_revision"`
	ToolVersion      string `json:"tool_version"`
}

type Settings interface {
	ListToolSettings(context.Context, Authority) ([]ToolSetting, error)
	UpdateToolSetting(context.Context, Authority, string, ToolSettingInput) (ToolSetting, error)
}

type SettingsBinding interface {
	Settings() Settings
	Availability() Availability
}
