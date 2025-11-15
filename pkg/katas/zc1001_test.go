package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/ast"
	"github.com/afadesigns/zshellcheck/pkg/lexer"
	"github.com/afadesigns/zshellcheck/pkg/parser"
)

func TestZC1001(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Violation
	}{
		{
			name:     "valid array access",
			input:    `echo ${my_array[1]}`,
			expected: []Violation{},
		},
		{
			name:  "invalid array access",
			input: `echo $my_array[1]`,
			expected: []Violation{
				{
					KataID:  "ZC1001",
					Message: "Use ${} for array element access. " +
						"Accessing array elements with `$my_array[1]` is not the correct syntax in Zsh.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			program := p.ParseProgram()

			violations := []Violation{}
			ast.Walk(program, func(node ast.Node) bool {
				if _, ok := node.(*ast.InvalidArrayAccess); ok {
					violations = append(violations, checkZC1001(node)...)
				}
				return true
			})

			if len(violations) != len(tt.expected) {
				t.Fatalf("expected %d violations, got %d", len(tt.expected), len(violations))
			}

			for i, v := range violations {
				if v.KataID != tt.expected[i].KataID {
					t.Errorf("expected kata ID %s, got %s", tt.expected[i].KataID, v.KataID)
				}
				if v.Message != tt.expected[i].Message {
					t.Errorf("expected message %s, got %s", tt.expected[i].Message, v.Message)
				}
				if v.Line != tt.expected[i].Line {
					t.Errorf("expected line %d, got %d", tt.expected[i].Line, v.Line)
				}
			}
		})
	}
}
