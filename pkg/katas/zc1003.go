package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1003",
		Title: "Use `((...))` for arithmetic comparisons instead of `[` or `test`",
		Description: "Bash/Zsh have a dedicated arithmetic context `((...))` which is cleaner and faster than `[` or `test` for numeric comparisons.",
		Check: checkZC1003,
	})
}

func checkZC1003(node ast.Node) []Violation {
	violations := []Violation{}

	if cmd, ok := node.(*ast.SimpleCommand); ok {
		if ident, ok := cmd.Name.(*ast.Identifier); ok {
			if ident.Value == "[" || ident.Value == "test" {
				for _, arg := range cmd.Arguments {
					val := getArgValue(arg)
					if val == "-eq" || val == "-ne" || val == "-lt" || val == "-le" || val == "-gt" || val == "-ge" {
						violations = append(violations, Violation{
							KataID:  "ZC1003",
							Message: "Use `((...))` for arithmetic comparisons instead of `[` or `test`.",
							Line:    ident.Token.Line,
							Column:  ident.Token.Column,
						})
						return violations
					}
				}
			}
		}
	}

	return violations
}
