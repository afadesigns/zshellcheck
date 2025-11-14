package ast

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/token"
)

func TestWalk(t *testing.T) {
	tests := []struct {
		name          string
		node          Node
		expectedCount int
	}{
		{
			"Let Statement",
			&Program{
				Statements: []Statement{
					&LetStatement{
						Token: token.Token{Type: token.LET, Literal: "let"},
						Name: &Identifier{
							Token: token.Token{Type: token.IDENT, Literal: "x"},
							Value: "x",
						},
						Value: &IntegerLiteral{
							Token: token.Token{Type: token.INT, Literal: "5"},
							Value: 5,
						},
					},
				},
			},
			4, // Program, LetStatement, Identifier, IntegerLiteral
		},
		{
			"If Statement",
			&IfStatement{
				Condition: &InfixExpression{
					Left:     &IntegerLiteral{Value: 1},
					Operator: "<",
					Right:    &IntegerLiteral{Value: 2},
				},
				Consequence: &BlockStatement{
					Statements: []Statement{
						&ExpressionStatement{
							Expression: &Identifier{Value: "a"},
						},
					},
				},
			},
			7, // If, Infix, Int, Int, Block, ExprStmt, Ident
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var nodes []Node
			Walk(tt.node, func(node Node) bool {
				nodes = append(nodes, node)
				return true
			})

			if len(nodes) != tt.expectedCount {
				t.Errorf("Walk did not visit all nodes. got=%d, want=%d", len(nodes), tt.expectedCount)
			}
		})
	}
}