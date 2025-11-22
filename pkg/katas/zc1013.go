package katas

import (

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.LetStatementNode, Kata{
		ID:          "ZC1013",
		Title:       "Use `((...))` for arithmetic operations instead of `let`",
		Description: "The `let` command is a shell builtin, but the `((...))` syntax is more portable " +
			"and generally preferred for arithmetic operations in Zsh.",
		Check:       checkZC1013,
	})
}

func checkZC1013(node ast.Node) []Violation {
	violations := []Violation{}

	if letStmt, ok := node.(*ast.LetStatement); ok {
		violations = append(violations, Violation{
			KataID:  "ZC1013",
			Message: "Use `((...))` for arithmetic operations instead of `let`.",
			Line:    letStmt.Token.Line,
			Column:  letStmt.Token.Column,
		})
	}

	return violations
}
