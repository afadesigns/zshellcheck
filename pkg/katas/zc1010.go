package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1010",
		Title: "Use `[[ ... ]]` instead of `[ ... ]` or `test`",
		Description: "The `[[` construct is more powerful and safer than `[` or `test`. " +
			"It prevents word splitting and glob expansion of its arguments.",
		Check: checkZC1010,
	})
}

func checkZC1010(node ast.Node) []Violation {
	violations := []Violation{}

	if cmd, ok := node.(*ast.SimpleCommand); ok {
		if ident, ok := cmd.Name.(*ast.Identifier); ok {
			if ident.Value == "[" || ident.Value == "test" {
				violations = append(violations, Violation{
					KataID:  "ZC1010",
					Message: "Use `[[ ... ]]` instead of `[ ... ]` or `test`. `[[` is safer and more powerful.",
					Line:    ident.Token.Line,
						Column:  ident.Token.Column,
				})
			}
		}
	}

	return violations
}
