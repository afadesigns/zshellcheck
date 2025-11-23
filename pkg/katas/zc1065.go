package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1065",
		Title: "Ensure spaces around `[` and `[[`",
		Description: "`[condition]` is parsed as a command named `[condition]`, which likely doesn't exist. " +
			"Add spaces: `[ condition ]`.",
		Check: checkZC1065,
	})
	// Register for DoubleBracketExpression to check `[[foo]]`
	RegisterKata(ast.DoubleBracketExpressionNode, Kata{
		ID:          "ZC1065",
		Title:       "Ensure spaces around `[` and `[[`",
		Description: "`[[condition]]` is parsed incorrectly. Add spaces: `[[ condition ]]`.",
		Check:       checkZC1065,
	})
}

func checkZC1065(node ast.Node) []Violation {
	violations := []Violation{}

	switch n := node.(type) {
	case *ast.SimpleCommand:
		if n.Name.String() == "[" {
			// Check first arg for preceding space
			if len(n.Arguments) > 0 {
				firstArg := n.Arguments[0]
				if !firstArg.TokenLiteralNode().HasPrecedingSpace {
					violations = append(violations, Violation{
						KataID:  "ZC1065",
						Message: "Missing space after `[`. Use `[ condition ]`.",
						Line:    n.Token.Line,
						Column:  n.Token.Column,
					})
				}
			}
		}
	case *ast.DoubleBracketExpression:
		// Check first expression
		if len(n.Expressions) > 0 {
			firstExp := n.Expressions[0]
			if !firstExp.TokenLiteralNode().HasPrecedingSpace {
				violations = append(violations, Violation{
					KataID:  "ZC1065",
					Message: "Missing space after `[[`. Use `[[ condition ]]`.",
					Line:    n.Token.Line,
					Column:  n.Token.Column,
				})
			}
		}
	}

	return violations
}
