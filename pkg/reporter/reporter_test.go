package reporter

import (
	"bytes"
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/ast"
	"reflect"
)

func TestTextReporter_Report(t *testing.T) {
	// Register a dummy kata for testing purposes
	katas.RegisterKata(reflect.TypeOf(&ast.Identifier{}), katas.Kata{
		ID:    "ZC9999",
		Title: "Test Kata",
	})

	violations := []katas.Violation{
		{
			KataID:  "ZC9999",
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