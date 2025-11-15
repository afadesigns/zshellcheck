package katas

import (
	"reflect"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(reflect.TypeOf(&ast.CommandSubstitution{}), Kata{
		ID:          "ZC1015",
		Title:       "Use `$(...)` for command substitution instead of backticks",
		Description: "The `$(...)` syntax is the modern, recommended way to perform command substitution. " +
			"It is more readable and can be nested easily, unlike backticks.",
		Check:       checkZC1015,
	})
}

func checkZC1015(node ast.Node) []Violation {
	violations := []Violation{}

	if cs, ok := node.(*ast.CommandSubstitution); ok {
		violations = append(violations, Violation{
			KataID:  "ZC1015",
			Message: "Use `$(...)` for command substitution instead of backticks.",
			Line:    cs.Token.Line,
			Column:  cs.Token.Column,
		})
	}

	return violations
}
