package toolsdk

import (
	"bytes"
	"testing"

	"github.com/domainry/domainry-tools-sdk/schema"
)

func TestScheduleDefinitionsAreStrictAndDoNotExposeOwnerRouting(t *testing.T) {
	definitions := ScheduleDefinitions()
	if len(definitions) != 7 {
		t.Fatalf("schedule definitions=%d", len(definitions))
	}
	seen := map[string]bool{}
	for _, definition := range definitions {
		if seen[definition.Key] || definition.Version != "2" || definition.ActionKey != "agent.conversation_tools."+definition.Key {
			t.Fatalf("invalid definition %+v", definition)
		}
		seen[definition.Key] = true
		if _, err := schema.CompileSchema(definition.InputSchema); err != nil {
			t.Fatalf("compile %s input: %v", definition.Key, err)
		}
		if _, err := schema.CompileSchema(definition.OutputSchema); err != nil {
			t.Fatalf("compile %s output: %v", definition.Key, err)
		}
		for _, forbidden := range [][]byte{[]byte(`"owner"`), []byte(`"workspace_id"`), []byte(`"user_id"`), []byte(`"target"`), []byte(`"allowed_actions"`), []byte(`"connection"`), []byte(`"provider"`), []byte(`"cron"`)} {
			if bytes.Contains(definition.InputSchema, forbidden) {
				t.Fatalf("%s input exposes %s", definition.Key, forbidden)
			}
		}
	}
	for _, key := range []string{ScheduleCreateToolKey, ScheduleListToolKey, ScheduleGetToolKey, ScheduleUpdateToolKey, SchedulePauseToolKey, ScheduleResumeToolKey, ScheduleDeleteToolKey} {
		if !seen[key] {
			t.Fatalf("missing %s", key)
		}
	}
}
