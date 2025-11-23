package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/ast"
	"github.com/afadesigns/zshellcheck/pkg/lexer"
	"github.com/afadesigns/zshellcheck/pkg/parser"
)

func TestCheckZC1037(t *testing.T) {
	tests := []struct {
		input    string
		expected []Violation
	}{
		{
			input: `echo $my_var`,
			expected: []Violation{
				{
					KataID:  "ZC1037",
					Message: "Unquoted variable expansion. Quote to prevent word splitting and globbing.",
					Line:    1,
					Column:  6,
				},
			},
		},
		{
			input:    `echo "$my_var"`,
			expected: []Violation{},
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()

		violations := []Violation{}
		ast.Walk(program, func(node ast.Node) bool {
			violations = append(violations, checkZC1037(node)...)
			return true
		})

		if len(violations) != len(tt.expected) {
			t.Fatalf("Expected %d violations, got %d for input: %s", len(tt.expected), len(violations), tt.input)
		}
	}
}
