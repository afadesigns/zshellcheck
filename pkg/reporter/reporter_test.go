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
	reporter := NewTextReporter(&buf)
	err := reporter.Report(violations)
	if err != nil {
		t.Fatalf("Report() returned an error: %v", err)
	}

	expected := "ZC9999: This is a test violation. (Test Kata)\n"
	if buf.String() != expected {
		t.Errorf("Report() produced incorrect output.\nGot: %q\nWant: %q", buf.String(), expected)
	}
}
