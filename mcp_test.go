package toolsdk

import (
	"bytes"
	"testing"

	"github.com/domainry/domainry-tools-sdk/schema"
)

func TestMCPDefinitionsAreFixedAndValid(t *testing.T) {
	definitions := MCPDefinitions()
	if len(definitions) != 3 {
		t.Fatalf("definitions=%d", len(definitions))
	}
	for _, definition := range definitions {
		if _, err := schema.CompileSchema(definition.InputSchema); err != nil {
			t.Fatalf("%s input schema: %v", definition.Key, err)
		}
		if _, err := schema.CompileSchema(definition.OutputSchema); err != nil {
			t.Fatalf("%s output schema: %v", definition.Key, err)
		}
		if bytes.Contains([]byte(definition.Description), []byte("remote description")) {
			t.Fatalf("remote data entered trusted definition %s", definition.Key)
		}
	}
	if definitions[0].Effect != "read" || definitions[1].Effect != "read" || definitions[2].Effect != "write" || definitions[2].Idempotency != "reconcile" {
		t.Fatalf("unexpected MCP effects: %#v", definitions)
	}
	definitions[0].InputSchema[0] = '['
	if bytes.Equal(definitions[0].InputSchema, MCPDefinitions()[0].InputSchema) {
		t.Fatal("definitions share mutable schema storage")
	}
}
