package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
	"github.com/afadesigns/zshellcheck/pkg/token"
)

func init() {
	RegisterKata(Kata{
		ID:          "ZC1001",
		Title:       "Prefer [[ over [ for Zsh-specific tests",
		Description: "The [[...]] construct is a Zsh keyword, offering safer and more powerful conditional expressions than the traditional [ command. It prevents word splitting and pathname expansion, and supports advanced features like regex matching.",
		Check:       checkZC1001,
	})
}

func checkZC1001(node ast.Node) []Violation {
	violations := []Violation{}

	ifExpression, ok := node.(*ast.IfExpression); ok {
		// We are looking for an IfExpression where the condition starts with a token.LBRACKET
		// This indicates the use of the old '[' command.
		if expressionStatement, ok := ifExpression.Condition.(*ast.ExpressionStatement); ok {
			if prefixExpression, ok := expressionStatement.Expression.(*ast.PrefixExpression); ok {
				if prefixExpression.Token.Type == token.LBRACKET {
					violations = append(violations, Violation{
						KataID:  "ZC1001",
						Message: "Prefer [[ over [ for Zsh-specific tests. [[ is a Zsh keyword, offering safer and more powerful conditional expressions.",
						Line:    prefixExpression.Token.Line,
						Column:  prefixExpression.Token.Column,
					})
				}
			}
		}
	}

	return violations
}
