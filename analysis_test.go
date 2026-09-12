package toolsdk

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestAnalysisInputSchemaIsClosedAndHasNoExecutableEscapeHatch(t *testing.T) {
	definition := AnalysisDefinitions()[0]
	var schema any
	if err := json.Unmarshal(definition.InputSchema, &schema); err != nil {
		t.Fatal(err)
	}
	forbidden := map[string]bool{"sql": true, "query": true, "code": true, "script": true, "javascript": true, "command": true}
	var inspect func(any, string)
	inspect = func(value any, path string) {
		switch current := value.(type) {
		case map[string]any:
			if properties, ok := current["properties"].(map[string]any); ok {
				if current["additionalProperties"] != false {
					t.Errorf("open object schema at %s", path)
				}
				for key := range properties {
					if forbidden[strings.ToLower(key)] {
						t.Errorf("executable input property %q at %s", key, path)
					}
				}
			}
			for key, child := range current {
				inspect(child, path+"/"+key)
			}
		case []any:
			for index, child := range current {
				inspect(child, path+"/"+string(rune('0'+index)))
			}
		}
	}
	inspect(schema, "$")
}
