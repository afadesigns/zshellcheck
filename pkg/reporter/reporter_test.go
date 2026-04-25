package reporter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/ast"
	"github.com/afadesigns/zshellcheck/pkg/config"
	"github.com/afadesigns/zshellcheck/pkg/katas"
)

func TestTextReporter_Report(t *testing.T) {
	// Register a dummy kata for testing purposes
	// Using a unique ID to avoid conflicts with existing katas if tests are run in parallel.
	const testKataID = "ZC9999"
	katas.RegisterKata(ast.IdentifierNode, katas.Kata{
		ID:    testKataID,
		Title: "Test Kata",
		Check: func(node ast.Node) []katas.Violation { return nil }, // Dummy check function
	})

	violations := []katas.Violation{
		{
			KataID:  testKataID,
			Message: "This is a test violation.",
			Level:   katas.SeverityWarning,
		},
	}

	var buf bytes.Buffer
	source := "first line\nsecond line\n"

	cfg := config.DefaultConfig()
	cfg.NoColor = false // Ensure colors are enabled for this test

	reporter := NewTextReporter(&buf, "test.zsh", source, cfg)
	// Update violations to include line/col
	violations[0].Line = 1
	violations[0].Column = 1

	err := reporter.Report(violations)
	if err != nil {
		t.Fatalf("Report() returned an error: %v", err)
	}

	// Update expected output to match new format with colors and context
	if !bytes.Contains(buf.Bytes(), []byte("This is a test violation.")) {
		t.Errorf("Report() produced incorrect output.\nGot:\n%s", buf.String())
	}

	// Check for code snippet and column pointer
	// Format:
	//   first line
	//   ↑
	expectedSnippet := "  first line\n"
	if !bytes.Contains(buf.Bytes(), []byte(expectedSnippet)) {
		t.Errorf("Report() output missing source code snippet.\nWant: %q\nGot:\n%q", expectedSnippet, buf.String())
	}

	// Check for column pointer with color (U+2191 upward arrow).
	expectedCaret := "  " + config.ColorBold + "↑" + config.ColorReset
	if !bytes.Contains(buf.Bytes(), []byte(expectedCaret)) {
		t.Errorf("Report() output missing caret.\nWant: %q\nGot:\n%q", expectedCaret, buf.String())
	}

	// Check for location
	expectedLocation := "test.zsh:1:1:"
	if !strings.Contains(buf.String(), expectedLocation) {
		t.Errorf("Report() output missing location.\nWant: %q\nGot:\n%q", expectedLocation, buf.String())
	}
}

func TestTextReporter_NoColor(t *testing.T) {
	violations := []katas.Violation{
		{
			KataID:  "ZC0001",
			Message: "no color test",
			Level:   katas.SeverityError,
			Line:    1,
			Column:  5,
		},
	}

	var buf bytes.Buffer
	cfg := config.DefaultConfig()
	cfg.NoColor = true

	reporter := NewTextReporter(&buf, "test.zsh", "echo hello", cfg)
	if err := reporter.Report(violations); err != nil {
		t.Fatalf("Report() error: %v", err)
	}

	output := buf.String()
	// With NoColor, there should be no ANSI escape sequences
	if strings.Contains(output, "\033[") {
		t.Errorf("expected no ANSI escape codes in NoColor mode, got:\n%s", output)
	}
	if !strings.Contains(output, "no color test") {
		t.Error("expected message in output")
	}
}

func TestTextReporter_AllSeverities(t *testing.T) {
	violations := []katas.Violation{
		{KataID: "ZC0001", Message: "error msg", Level: katas.SeverityError, Line: 1, Column: 1},
		{KataID: "ZC0002", Message: "warning msg", Level: katas.SeverityWarning, Line: 1, Column: 1},
		{KataID: "ZC0003", Message: "info msg", Level: katas.SeverityInfo, Line: 1, Column: 1},
	}

	var buf bytes.Buffer
	cfg := config.DefaultConfig()
	reporter := NewTextReporter(&buf, "test.zsh", "echo hello", cfg)
	if err := reporter.Report(violations); err != nil {
		t.Fatalf("Report() error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "error msg") {
		t.Error("missing error message")
	}
	if !strings.Contains(output, "warning msg") {
		t.Error("missing warning message")
	}
	if !strings.Contains(output, "info msg") {
		t.Error("missing info message")
	}
}

func TestTextReporter_EmptyViolations(t *testing.T) {
	var buf bytes.Buffer
	cfg := config.DefaultConfig()
	reporter := NewTextReporter(&buf, "test.zsh", "echo hello", cfg)
	if err := reporter.Report(nil); err != nil {
		t.Fatalf("Report() error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty output for no violations, got: %q", buf.String())
	}
}

func TestTextReporter_LineOutOfRange(t *testing.T) {
	violations := []katas.Violation{
		{KataID: "ZC0001", Message: "out of range", Level: katas.SeverityWarning, Line: 99, Column: 1},
	}

	var buf bytes.Buffer
	cfg := config.DefaultConfig()
	reporter := NewTextReporter(&buf, "test.zsh", "one line", cfg)
	if err := reporter.Report(violations); err != nil {
		t.Fatalf("Report() error: %v", err)
	}

	// Should still report the message, just no code snippet
	output := buf.String()
	if !strings.Contains(output, "out of range") {
		t.Error("expected message in output even with out-of-range line")
	}
}

func TestTextReporter_ZeroColumn(t *testing.T) {
	violations := []katas.Violation{
		{KataID: "ZC0001", Message: "zero col", Level: katas.SeverityWarning, Line: 1, Column: 0},
	}

	var buf bytes.Buffer
	cfg := config.DefaultConfig()
	reporter := NewTextReporter(&buf, "test.zsh", "echo hello", cfg)
	if err := reporter.Report(violations); err != nil {
		t.Fatalf("Report() error: %v", err)
	}

	// Column 0 should produce padding of 0 (caret at start)
	if !strings.Contains(buf.String(), "zero col") {
		t.Error("expected message in output")
	}
}

func TestTextReporter_MultiLineSource(t *testing.T) {
	violations := []katas.Violation{
		{KataID: "ZC0001", Message: "line2 issue", Level: katas.SeverityError, Line: 2, Column: 3},
	}

	var buf bytes.Buffer
	cfg := config.DefaultConfig()
	reporter := NewTextReporter(&buf, "script.zsh", "line one\nline two\nline three", cfg)
	if err := reporter.Report(violations); err != nil {
		t.Fatalf("Report() error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "line two") {
		t.Error("expected second line snippet in output")
	}
	if !strings.Contains(output, "script.zsh:2:3:") {
		t.Error("expected correct location in output")
	}
}

func TestSarifReporter_Report(t *testing.T) {
	violations := []katas.Violation{
		{KataID: "ZC1001", Message: "test violation 1", Level: katas.SeverityError, Line: 1, Column: 1},
		{KataID: "ZC1002", Message: "test violation 2", Level: katas.SeverityWarning, Line: 5, Column: 10},
	}

	var buf bytes.Buffer
	reporter := NewSarifReporter(&buf, "test.zsh")
	if err := reporter.Report(violations); err != nil {
		t.Fatalf("Report() error: %v", err)
	}

	// Verify it's valid JSON
	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	// Check SARIF version
	if version, ok := result["version"].(string); !ok || version != "2.1.0" {
		t.Errorf("expected SARIF version 2.1.0, got %v", result["version"])
	}

	// Check runs array
	runs, ok := result["runs"].([]interface{})
	if !ok || len(runs) != 1 {
		t.Fatalf("expected 1 run, got %v", result["runs"])
	}

	run := runs[0].(map[string]interface{})
	// Check tool name
	tool := run["tool"].(map[string]interface{})
	driver := tool["driver"].(map[string]interface{})
	if driver["name"] != "zshellcheck" {
		t.Errorf("expected tool name zshellcheck, got %v", driver["name"])
	}

	// Check results count
	results := run["results"].([]interface{})
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}

	// Verify first result
	r0 := results[0].(map[string]interface{})
	if r0["ruleId"] != "ZC1001" {
		t.Errorf("expected ruleId ZC1001, got %v", r0["ruleId"])
	}
	if r0["message"] != "test violation 1" {
		t.Errorf("expected message 'test violation 1', got %v", r0["message"])
	}
}

func TestSarifReporter_EmptyViolations(t *testing.T) {
	var buf bytes.Buffer
	reporter := NewSarifReporter(&buf, "test.zsh")
	if err := reporter.Report(nil); err != nil {
		t.Fatalf("Report() error: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	runs := result["runs"].([]interface{})
	run := runs[0].(map[string]interface{})
	// Empty violations should produce empty results array (not nil in JSON)
	results := run["results"].([]interface{})
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestJSONReporter_EmptyViolations(t *testing.T) {
	var buf bytes.Buffer
	reporter := NewJSONReporter(&buf)
	if err := reporter.Report([]katas.Violation{}); err != nil {
		t.Fatalf("Report() error: %v", err)
	}

	var result []interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("expected empty array, got %d items", len(result))
	}
}

type failWriter struct{}

func (w *failWriter) Write(p []byte) (n int, err error) {
	return 0, fmt.Errorf("write failed")
}

func TestTextReporter_WriterError(t *testing.T) {
	violations := []katas.Violation{
		{KataID: "ZC0001", Message: "test", Level: katas.SeverityError, Line: 1, Column: 1},
	}

	cfg := config.DefaultConfig()
	reporter := NewTextReporter(&failWriter{}, "test.zsh", "echo hello", cfg)
	err := reporter.Report(violations)
	if err == nil {
		t.Error("expected error from failing writer")
	}
}

func TestJSONReporter_WriterError(t *testing.T) {
	violations := []katas.Violation{
		{KataID: "ZC0001", Message: "test", Level: katas.SeverityError, Line: 1, Column: 1},
	}

	reporter := NewJSONReporter(&failWriter{})
	err := reporter.Report(violations)
	if err == nil {
		t.Error("expected error from failing writer")
	}
}

func TestSarifReporter_WriterError(t *testing.T) {
	violations := []katas.Violation{
		{KataID: "ZC0001", Message: "test", Level: katas.SeverityError, Line: 1, Column: 1},
	}

	reporter := NewSarifReporter(&failWriter{}, "test.zsh")
	err := reporter.Report(violations)
	if err == nil {
		t.Error("expected error from failing writer")
	}
}

func TestJSONReporter_MultipleViolations(t *testing.T) {
	violations := []katas.Violation{
		{KataID: "ZC0001", Message: "first", Level: katas.SeverityError, Line: 1, Column: 1},
		{KataID: "ZC0002", Message: "second", Level: katas.SeverityWarning, Line: 2, Column: 5},
		{KataID: "ZC0003", Message: "third", Level: katas.SeverityInfo, Line: 3, Column: 10},
	}

	var buf bytes.Buffer
	reporter := NewJSONReporter(&buf)
	if err := reporter.Report(violations); err != nil {
		t.Fatalf("Report() error: %v", err)
	}

	var result []katas.Violation
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	if len(result) != 3 {
		t.Fatalf("expected 3 violations, got %d", len(result))
	}
	if result[0].KataID != "ZC0001" {
		t.Errorf("expected first KataID=ZC0001, got %s", result[0].KataID)
	}
	if result[2].Message != "third" {
		t.Errorf("expected third message='third', got %s", result[2].Message)
	}
}
