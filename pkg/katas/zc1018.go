package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1018",
		Title: "Use `((...))` for C-style arithmetic instead of `expr`",
		Description: "The `((...))` construct in Zsh allows for C-style arithmetic. " +
			"It is generally more efficient and readable than using `expr` or other " +
			"external commands for arithmetic.",
		Check: checkZC1018,
	})
}

func checkZC1018(node ast.Node) []Violation {
	violations := []Violation{}

	if cmd, ok := node.(*ast.SimpleCommand); ok {
		if ident, ok := cmd.Name.(*ast.Identifier); ok {
			if ident.Value == "expr" {
				violations = append(violations, Violation{
					KataID:  "ZC1018",
					Message: "Use `((...))` for C-style arithmetic instead of `expr`.",
					Line:    ident.Token.Line,
					Column:  ident.Token.Column,
				})
			}
		}
	}

	return violations
}
