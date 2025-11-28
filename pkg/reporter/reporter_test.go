package reporter

import (
	"bytes"
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/ast"
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
		},
	}

	var buf bytes.Buffer
	source := "first line\nsecond line\n"
	reporter := NewTextReporter(&buf, "test.zsh", source)
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
	if !bytes.Contains(buf.Bytes(), []byte("first line")) {
		t.Errorf("Report() output missing source line.\nGot:\n%s", buf.String())
	}
}
