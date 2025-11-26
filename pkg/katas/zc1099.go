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
		if _, ok := infix.Right.(*ast.WhileLoopStatement); ok {
			// Check if the while loop condition is `read ...`
			// while loop struct has `Condition *BlockStatement`
			// `read` is a SimpleCommand.
			// We need to peek into the condition block.
			return []Violation{{
				KataID:  "ZC1099",
				Message: "Consider using `for line in ${(f)variable}` instead of `... | while read line`. It's faster and cleaner in Zsh.",
				Line:    infix.Token.Line,
				Column:  infix.Token.Column,
			}}
		}
	}
	return nil
}
