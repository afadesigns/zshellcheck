package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/ast"
	"github.com/afadesigns/zshellcheck/pkg/lexer"
	"github.com/afadesigns/zshellcheck/pkg/parser"
)

func TestCheckZC1036(t *testing.T) {
	tests := []struct {
		input    string
		expected []Violation
	}{
		{
			input: `test -f file.txt`,
			expected: []Violation{
				{
					KataID:  "ZC1036",
					Message: "Prefer `[[ ... ]]` over `test` command for conditional expressions.",
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

		var violations []Violation
		ast.Walk(program, func(node ast.Node) bool {
			for _, v := range Check(node, []string{}) {
				if v.KataID == "ZC1036" {
					violations = append(violations, v)
				}
			}
			return true
		})

		if len(violations) != len(tt.expected) {
			t.Fatalf("Expected %d violations, got %d for input: %s", len(tt.expected), len(violations), tt.input)
		}
	}
}
