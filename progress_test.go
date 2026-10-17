package toolsdk_test

import (
	"context"
	"encoding/json"
	"testing"

	toolsdk "github.com/domainry/domainry-tools-sdk"
)

func TestPublishedToolCanDeclareAndReportConditionalProgress(t *testing.T) {
	definition := toolsdk.Definition{ExecutionPlan: []toolsdk.ProgressStep{
		{Key: "read", Label: "Read source"},
		{Key: "enrich", Label: "Enrich result", Optional: true},
	}}
	raw, err := json.Marshal(definition)
	if err != nil || !json.Valid(raw) || len(definition.ExecutionPlan) != 2 || !definition.ExecutionPlan[1].Optional {
		t.Fatalf("definition=%#v raw=%s err=%v", definition, raw, err)
	}

	var received []toolsdk.Progress
	ctx := toolsdk.WithProgressReporter(context.Background(), func(_ context.Context, progress toolsdk.Progress) error {
		received = append(received, progress)
		return nil
	})
	for _, progress := range []toolsdk.Progress{
		{Key: "read", Status: toolsdk.ProgressStatusStarted},
		{Key: "read", Status: toolsdk.ProgressStatusCompleted},
		{Key: "enrich", Status: toolsdk.ProgressStatusSkipped},
	} {
		if err = toolsdk.ReportProgress(ctx, progress); err != nil {
			t.Fatal(err)
		}
	}
	if len(received) != 3 || received[2].Status != toolsdk.ProgressStatusSkipped {
		t.Fatalf("received=%#v", received)
	}
}
