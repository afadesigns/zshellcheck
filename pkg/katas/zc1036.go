package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1036",
		Title: "Prefer `[[ ... ]]` over `test` command",
		Description: "The `[[ ... ]]` construct is a more powerful and safer alternative to the `test` " +
			"command (or `[ ... ]`) for conditional expressions in modern shells. It handles word " +
			"splitting and globbing more intuitively and supports advanced features like regex matching.",
		Check: checkZC1036,
	})
}

func checkZC1036(node ast.Node) []Violation {
	violations := []Violation{}

	if cmd, ok := node.(*ast.SimpleCommand); ok {
		if cmd.Name.TokenLiteral() == "test" {
			violations = append(violations, Violation{
				KataID:  "ZC1036",
				Message: "Prefer `[[ ... ]]` over `test` command for conditional expressions.",
				Line:    cmd.Token.Line,
				Column:  cmd.Token.Column,
			})
		}
	}

	return violations
}
