// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package reporter

import (
	"bytes"
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
