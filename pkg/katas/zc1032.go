package katas

import (

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.LetStatementNode, Kata{
		ID:          "ZC1032",
		Title:       "Use `((...))` for C-style incrementing",
		Description: "Instead of `let i=i+1`, you can use the more concise and idiomatic C-style " +
			"increment `(( i++ ))` in Zsh.",
		Check:       checkZC1032,
	})
}

func checkZC1032(node ast.Node) []Violation {
	violations := []Violation{}

	if letStmt, ok := node.(*ast.LetStatement); ok {
		if infixExpr, ok := letStmt.Value.(*ast.InfixExpression); ok {
			if leftIdent, ok := infixExpr.Left.(*ast.Identifier); ok {
				if rightInt, ok := infixExpr.Right.(*ast.IntegerLiteral); ok {
					if letStmt.Name.Value == leftIdent.Value && infixExpr.Operator == "+" && rightInt.Value == 1 {
						violations = append(violations, Violation{
							KataID:  "ZC1032",
							Message: "Use `(( i++ ))` for C-style incrementing instead of `let i=i+1`.",
							Line:    letStmt.Token.Line,
							Column:  letStmt.Token.Column,
						})
					}
				}
			}
		}
	}

	return violations
}
