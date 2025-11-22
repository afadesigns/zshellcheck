package katas

import (

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1006",
		Title: "Prefer [[ over test for tests",
		Description: "The `test` command is an external command and may not be available on all systems. " +
			"The `[[...]]` construct is a Zsh keyword, offering safer and more powerful conditional " +
			"expressions than the traditional `test` command. It prevents word splitting and pathname " +
			"expansion, and supports advanced features like regex matching.",
		Check: checkZC1006,
	})
}

func checkZC1006(node ast.Node) []Violation {
	violations := []Violation{}

	if cmd, ok := node.(*ast.SimpleCommand); ok {
		if ident, ok := cmd.Name.(*ast.Identifier); ok {
			if ident.Value == "test" {
				violations = append(violations, Violation{
					KataID: "ZC1006",
					Message: "Prefer [[ over test for tests. " +
						"[[ is a Zsh keyword that offers safer and more powerful conditional expressions.",
					Line:   ident.Token.Line,
					Column: ident.Token.Column,
				})
			}
		}
	}

	return violations
}
