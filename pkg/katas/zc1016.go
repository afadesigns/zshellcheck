package katas

import (
	"reflect"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(reflect.TypeOf(&ast.BracketExpression{}), Kata{
		ID:          "ZC1016",
		Title:       "Use `[[ ... ]]` for tests instead of `[ ... ]`",
		Description: "The `[[ ... ]]` construct is a Zsh keyword and is generally safer and more powerful " +
			"than the `[ ... ]` command. It prevents word splitting and pathname expansion, and supports " +
			"advanced features like regex matching.",
		Check:       checkZC1016,
	})
}

func checkZC1016(node ast.Node) []Violation {
	violations := []Violation{}

	if be, ok := node.(*ast.BracketExpression); ok {
		violations = append(violations, Violation{
			KataID:  "ZC1016",
			Message: "Use `[[ ... ]]` for tests instead of `[ ... ]`.",
			Line:    be.Token.Line,
			Column:  be.Token.Column,
		})
	}

	return violations
}
