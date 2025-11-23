package testutil

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/ast"
	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/lexer"
	"github.com/afadesigns/zshellcheck/pkg/parser"
)

func Check(code string, kataID string) []katas.Violation {
	l := lexer.New(code)
	p := parser.New(l)
	program := p.ParseProgram()

	var violations []katas.Violation
	ast.Walk(program, func(node ast.Node) bool {
		if node == nil {
			return true
		}
		if katasForNode, ok := katas.Registry.KatasByNodeType[node.Type()]; ok {
			for _, kata := range katasForNode {
				if kata.ID == kataID {
					violations = append(violations, kata.Check(node)...)
				}
			}
		}
		return true
	})

	var result []katas.Violation
	for _, v := range violations {
		result = append(result, katas.Violation{
			KataID:  v.KataID,
			Message: v.Message,
			Line:    v.Line,
			Column:  v.Column,
		})
	}
	return result
}

func AssertViolations(t *testing.T, code string, actual []katas.Violation, expected []katas.Violation) {
	t.Helper()

	if len(actual) != len(expected) {
		t.Fatalf("expected %d violations, got %d", len(expected), len(actual))
	}

	for i, v := range actual {
		if v.KataID != expected[i].KataID {
			t.Errorf("expected kata ID %q, got %q", expected[i].KataID, v.KataID)
		}
		if v.Message != expected[i].Message {
			t.Errorf("expected message %q, got %q", expected[i].Message, v.Message)
		}
		if v.Line != expected[i].Line {
			t.Errorf("expected line %d, got %d", expected[i].Line, v.Line)
		}
		if v.Column != expected[i].Column {
			t.Errorf("expected column %d, got %d", expected[i].Column, v.Column)
		}
	}
}
