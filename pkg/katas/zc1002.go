package katas

import (
	"reflect"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(reflect.TypeOf(&ast.CommandSubstitution{}), Kata{
		ID:    "ZC1002",
		Title: "Use $(...) instead of backticks for command substitution",
		Description: "The `$(...)` syntax is the modern, recommended way to perform command substitution. " +
			"It is more readable and can be nested easily, unlike backticks.",
		Check: checkZC1002,
	})
}

func checkZC1002(node ast.Node) []Violation {
	violations := []Violation{}

	if commandSubstitution, ok := node.(*ast.CommandSubstitution); ok {
		violations = append(violations, Violation{
			KataID: "ZC1002",
			Message: "Use $(...) instead of backticks for command substitution. " +
				"The `$(...)` syntax is more readable and can be nested easily.",
			Line:   commandSubstitution.Token.Line,
			Column: commandSubstitution.Token.Column,
		})
	}

	return violations
}
