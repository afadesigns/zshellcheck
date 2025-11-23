package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.LetStatementNode, Kata{
		ID:    "ZC1026",
		Title: "Use `$((...))` for arithmetic expansion",
		Description: "The `$((...))` syntax is the modern, recommended way to perform arithmetic expansion. " +
			"It is more readable and can be nested easily, unlike `let`.",
		Check: checkZC1026,
	})
}

func checkZC1026(node ast.Node) []Violation {
	violations := []Violation{}

	if let, ok := node.(*ast.LetStatement); ok {
		violations = append(violations, Violation{
			KataID:  "ZC1026",
			Message: "Use `$((...))` for arithmetic expansion instead of `let`.",
			Line:    let.Token.Line,
			Column:  let.Token.Column,
		})
	}

	return violations
}
