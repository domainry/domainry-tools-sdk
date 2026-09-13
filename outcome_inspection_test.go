package toolsdk

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestOutcomeInspectionGuardCannotCrossUnrelatedJSONTransport(t *testing.T) {
	raw, err := json.Marshal(Request{OutcomeInspectionToken: "server-inspection-token", IdempotencyKey: "original-operation"})
	if err != nil || strings.Contains(string(raw), "inspection") {
		t.Fatalf("local guard escaped: %s %v", raw, err)
	}
	var in Request
	if err = json.Unmarshal([]byte(`{"OutcomeInspectionToken":"caller-token","IdempotencyKey":"original-operation"}`), &in); err != nil {
		t.Fatal(err)
	}
	if in.OutcomeInspectionToken != "" || in.IdempotencyKey != "original-operation" {
		t.Fatalf("client supplied local guard: %+v", in)
	}
}
