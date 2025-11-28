package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.InfixExpressionNode, Kata{
		ID:    "ZC1099",
		Title: "Use `(f)` flag to split lines instead of `while read`",
		Description: "Zsh provides the `(f)` parameter expansion flag to split a string into lines. " +
			"Iterating over `${(f)variable}` is often cleaner and faster than piping to `while read`.",
		Check: checkZC1099,
	})
}

func checkZC1099(node ast.Node) []Violation {
	infix, ok := node.(*ast.InfixExpression)
	if !ok {
		return nil
	}

	if infix.Operator == "|" {
		if whileLoop, ok := infix.Right.(*ast.WhileLoopStatement); ok {
			foundRead := false
			for _, stmt := range whileLoop.Condition.(*ast.BlockStatement).Statements {
				if exprStmt, ok := stmt.(*ast.ExpressionStatement); ok {
					if simpleCmd, ok := exprStmt.Expression.(*ast.SimpleCommand); ok && simpleCmd.Name != nil && simpleCmd.Name.String() == "read" {
						foundRead = true
						break
					}
				}
			}

			if foundRead {
				return []Violation{{
					KataID:  "ZC1099",
					Message: "Consider using `for line in ${(f)variable}` instead of `... | while read line`. It's faster and cleaner in Zsh.",
					Line:    infix.Token.Line,
					Column:  infix.Token.Column,
				}}
			}
		}
	}
	return nil
}
