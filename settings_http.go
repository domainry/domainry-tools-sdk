package toolsdk

import (
	action "github.com/domainry/domainry-foundation/action"
	"github.com/domainry/domainry-foundation/modulehttp"
)

const ActionToolSettingsList = "tools.preferences.list"
const ActionToolSettingsUpdate = "tools.preferences.update"

func ToolSettingsRoutes() []modulehttp.Route {
	out := []modulehttp.Route{}
	for _, spec := range []struct{ key, operation, method, path, label string }{{ActionToolSettingsList, "list", "GET", "/tools/preferences", "查看自己的工具设置"}, {ActionToolSettingsUpdate, "update", "PUT", "/tools/preferences/{toolKey}", "修改自己的工具设置"}} {
		effect := action.EffectRead
		idempotency := "not_applicable"
		if spec.method == "PUT" {
			effect = action.EffectWrite
			idempotency = "expected_preference_revision"
		}
		out = append(out, modulehttp.Route{Action: action.ActionDefinition{
			Key: spec.key, Owner: "module:tools", SourceKind: "module_http", CapabilityKey: "tools.preferences", CapabilityLabel: "工具设置", OperationKey: spec.operation, OperationLabel: spec.label, Label: spec.label,
			Exposures: []action.Exposure{action.ExposurePublic}, HTTP: &action.HTTPBinding{Method: spec.method, RouteTemplate: spec.path}, EffectClass: effect, RiskLevel: action.RiskLow, IdempotencyDecision: idempotency, AuditClass: "tools_owner_preference", LifecycleStatus: action.LifecycleActive,
			Authorization: action.Authorization{Strategy: action.AuthorizationAuthenticated}, Permission: &action.PermissionDefinition{Key: spec.key, Owner: "module:tools", ResourceKey: "tools.preferences", OperationKey: spec.operation, Label: spec.label, Category: "Tools", LifecycleStatus: action.LifecycleActive},
		}})
	}
	return out
}
