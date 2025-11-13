package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/ast"
	"github.com/afadesigns/zshellcheck/pkg/lexer"
	"github.com/afadesigns/zshellcheck/pkg/parser"
)

func TestCheckZC1012(t *testing.T) {
	tests := []struct {
		input    string
		expected []Violation
	}{
		{
			input: "`ls`",
			expected: []Violation{
				{
					KataID:  "ZC1012",
					Message: "Use `$(command)` instead of backticks for command substitution.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()

		violations := []Violation{}
		ast.Walk(program, func(node ast.Node) bool {
			violations = append(violations, checkZC1012(node)...)
			return true
		})

		if len(violations) != len(tt.expected) {
			t.Fatalf("Expected %d violations, got %d for input: %s", len(tt.expected), len(violations), tt.input)
		}

		for i, v := range violations {
			if v.KataID != tt.expected[i].KataID {
				t.Errorf("Expected KataID %s, got %s", tt.expected[i].KataID, v.KataID)
			}
			if v.Message != tt.expected[i].Message {
				t.Errorf("Expected Message %s, got %s", tt.expected[i].Message, v.Message)
			}
			if v.Line != tt.expected[i].Line {
				t.Errorf("Expected Line %d, got %d", tt.expected[i].Line, v.Line)
			}
			if v.Column != tt.expected[i].Column {
				t.Errorf("Expected Column %d, got %d", tt.expected[i].Column, v.Column)
			}
		}
	}
}
