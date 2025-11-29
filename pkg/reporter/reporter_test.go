package reporter

import (
	"bytes"
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
			Level:   katas.Warning,
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
	
	// Check for code snippet and caret
	// Format:
	//   first line
	//   ^
	expectedSnippet := "  first line\n"
	if !bytes.Contains(buf.Bytes(), []byte(expectedSnippet)) {
		t.Errorf("Report() output missing source code snippet.\nWant: %q\nGot:\n%q", expectedSnippet, buf.String())
	}
	
	// Check for caret with color
	// Note: config constants are not exported, so using hardcoded ANSI codes or config public constants if available.
	// pkg/config/config.go exports ColorBold etc.
	expectedCaret := "  " + config.ColorBold + "^" + config.ColorReset
	if !bytes.Contains(buf.Bytes(), []byte(expectedCaret)) {
		t.Errorf("Report() output missing caret.\nWant: %q\nGot:\n%q", expectedCaret, buf.String())
	}
    
    // Check for location
    expectedLocation := "test.zsh:1:1:"
    if !strings.Contains(buf.String(), expectedLocation) {
         t.Errorf("Report() output missing location.\nWant: %q\nGot:\n%q", expectedLocation, buf.String())
    }
}
