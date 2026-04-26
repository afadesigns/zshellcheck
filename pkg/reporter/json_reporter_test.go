// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package reporter

import (
	"bytes"
	"encoding/json"
	"strings"
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

func TestTextReporter_StyleSeverity(t *testing.T) {
	violations := []katas.Violation{
		{KataID: "ZC0001", Message: "style msg", Level: katas.SeverityStyle, Line: 1, Column: 1},
	}

	var buf bytes.Buffer
	cfg := config.DefaultConfig()
	reporter := NewTextReporter(&buf, "test.zsh", "echo hello", cfg)
	if err := reporter.Report(violations); err != nil {
		t.Fatalf("Report() error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "style msg") {
		t.Error("expected style message in output")
	}
	// Style severity uses cyan color
	if !strings.Contains(output, "\033[36m") {
		t.Error("expected cyan color code for style severity")
	}
}
