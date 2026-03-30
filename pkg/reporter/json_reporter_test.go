package reporter

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/config"
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

func TestTextReporter_Report_AllSeverities(t *testing.T) {
	source := "echo hello\necho world\necho done"
	violations := []katas.Violation{
		{KataID: "ZC0001", Message: "error msg", Line: 1, Column: 1, Level: katas.SeverityError},
		{KataID: "ZC0002", Message: "warn msg", Line: 2, Column: 6, Level: katas.SeverityWarning},
		{KataID: "ZC0003", Message: "info msg", Line: 3, Column: 1, Level: katas.SeverityInfo},
		{KataID: "ZC0004", Message: "style msg", Line: 1, Column: 5, Level: katas.SeverityStyle},
	}

	var buf bytes.Buffer
	cfg := config.DefaultConfig()
	reporter := NewTextReporter(&buf, "test.zsh", source, cfg)
	err := reporter.Report(violations)
	if err != nil {
		t.Fatalf("Report() returned an error: %v", err)
	}

	output := buf.String()

	// Check that all violations appear in the output
	for _, v := range violations {
		if !bytes.Contains([]byte(output), []byte(v.KataID)) {
			t.Errorf("output missing KataID %q", v.KataID)
		}
		if !bytes.Contains([]byte(output), []byte(v.Message)) {
			t.Errorf("output missing message %q", v.Message)
		}
	}
}

func TestTextReporter_Report_NoColor(t *testing.T) {
	source := "echo hello"
	violations := []katas.Violation{
		{KataID: "ZC0001", Message: "test", Line: 1, Column: 1, Level: katas.SeverityError},
	}

	var buf bytes.Buffer
	cfg := config.DefaultConfig()
	cfg.NoColor = true
	reporter := NewTextReporter(&buf, "test.zsh", source, cfg)
	err := reporter.Report(violations)
	if err != nil {
		t.Fatalf("Report() returned an error: %v", err)
	}

	output := buf.String()
	// Should NOT contain ANSI escape codes
	if bytes.Contains([]byte(output), []byte("\033[")) {
		t.Error("output contains ANSI escape codes when NoColor is true")
	}
}

func TestTextReporter_Report_EmptyViolations(t *testing.T) {
	var buf bytes.Buffer
	cfg := config.DefaultConfig()
	reporter := NewTextReporter(&buf, "test.zsh", "echo hello", cfg)
	err := reporter.Report([]katas.Violation{})
	if err != nil {
		t.Fatalf("Report() returned an error: %v", err)
	}
	if buf.Len() != 0 {
		t.Error("expected empty output for no violations")
	}
}

func TestTextReporter_Report_InvalidLine(t *testing.T) {
	// Line number beyond the source lines -- should not panic
	source := "echo hello"
	violations := []katas.Violation{
		{KataID: "ZC0001", Message: "test", Line: 99, Column: 1, Level: katas.SeverityWarning},
	}

	var buf bytes.Buffer
	cfg := config.DefaultConfig()
	reporter := NewTextReporter(&buf, "test.zsh", source, cfg)
	err := reporter.Report(violations)
	if err != nil {
		t.Fatalf("Report() returned an error: %v", err)
	}
}

func TestTextReporter_Report_ZeroColumn(t *testing.T) {
	// Column 0 should not cause negative padding
	source := "echo hello"
	violations := []katas.Violation{
		{KataID: "ZC0001", Message: "test", Line: 1, Column: 0, Level: katas.SeverityInfo},
	}

	var buf bytes.Buffer
	cfg := config.DefaultConfig()
	reporter := NewTextReporter(&buf, "test.zsh", source, cfg)
	err := reporter.Report(violations)
	if err != nil {
		t.Fatalf("Report() returned an error: %v", err)
	}
}

func TestTextReporter_Report_UnknownSeverity(t *testing.T) {
	// Unknown severity level uses no color
	source := "echo hello"
	violations := []katas.Violation{
		{KataID: "ZC0001", Message: "test", Line: 1, Column: 1, Level: "unknown"},
	}

	var buf bytes.Buffer
	cfg := config.DefaultConfig()
	reporter := NewTextReporter(&buf, "test.zsh", source, cfg)
	err := reporter.Report(violations)
	if err != nil {
		t.Fatalf("Report() returned an error: %v", err)
	}
}
