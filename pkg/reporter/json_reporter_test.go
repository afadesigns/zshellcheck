package reporter

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
)

func TestJSONReporter_Report(t *testing.T) {
	violations := []katas.Violation{
		{
			KataID:  "ZC1001",
			Message: "This is a test violation.",
			Line:    1,
			Column:  1,
		},
	}

	var buf bytes.Buffer
	reporter := NewJSONReporter(&buf)
	err := reporter.Report(violations)
	if err != nil {
		t.Fatalf("Report() returned an error: %v", err)
	}

	var reportedViolations []katas.Violation
	err = json.Unmarshal(buf.Bytes(), &reportedViolations)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if len(reportedViolations) != 1 {
		t.Fatalf("Expected 1 violation, got %d", len(reportedViolations))
	}

	if reportedViolations[0].KataID != "ZC1001" {
		t.Errorf("Expected KataID to be ZC1001, got %s", reportedViolations[0].KataID)
	}
}
