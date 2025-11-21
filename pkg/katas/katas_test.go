package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/ast"
	"github.com/afadesigns/zshellcheck/pkg/lexer"
	"github.com/afadesigns/zshellcheck/pkg/parser"
)

// check is the single, authoritative test helper. It parses a string,
// walks the resulting AST, and runs the non-recursive Check function on
// every node, returning all found violations. This perfectly mimics the
// application's main execution loop.
func check(input string, registry *KatasRegistry, kataID string) []Violation {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	var violations []Violation
	ast.Walk(program, func(node ast.Node) bool {
		foundViolations := registry.Check(node, []string{})
		for _, v := range foundViolations {
			if v.KataID == kataID {
				violations = append(violations, v)
			}
		}
		return true
	})
	return violations
}

// assertViolations checks that the violations found by the linter match
// the expected violations.
func assertViolations(t *testing.T, input string, violations []Violation, expected []Violation) {
	t.Helper()

	if len(violations) != len(expected) {
		l := lexer.New(input)
		p := parser.New(l)
		program := p.ParseProgram()
		t.Fatalf("Expected %d violations, but got %d for input:\n%s\nAST:\n%s",
			len(expected), len(violations), input, program.String())
	}

	for i, v := range violations {
		if v.KataID != expected[i].KataID {
			t.Errorf("Violation %d: Expected KataID %s, got %s", i, expected[i].KataID, v.KataID)
		}
		if v.Message != expected[i].Message {
			t.Errorf("Violation %d: Expected Message %q, got %q", i, expected[i].Message, v.Message)
		}
		if v.Line != expected[i].Line {
			t.Errorf("Violation %d: Expected Line %d, got %d", i, expected[i].Line, v.Line)
		}
		if v.Column != expected[i].Column {
			t.Errorf("Violation %d: Expected Column %d, got %d", i, expected[i].Column, v.Column)
		}
	}
}
