package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:          "ZC1010",
		Title:       "Use [[ ... ]] instead of [ ... ]",
		Description: "Zsh's [[ ... ]] is more powerful and safer than [ ... ]. It supports pattern matching, regex, and doesn't require quoting variables to prevent word splitting.",
		Check:       checkZC1010,
	})
}

func checkZC1010(node ast.Node) []Violation {
	violations := []Violation{}

	if cmd, ok := node.(*ast.SimpleCommand); ok {
        // Check if command name is "["
        if cmd.Name.String() == "[" {
            violations = append(violations, Violation{
                KataID:  "ZC1010",
                Message: "Use `[[ ... ]]` instead of `[ ... ]` or `test`. `[[` is safer and more powerful.",
                Line:    cmd.Token.Line,
                Column:  cmd.Token.Column,
            })
        }
	}

	return violations
}
