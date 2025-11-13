package katas

import (
	"reflect"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(reflect.TypeOf(&ast.CommandSubstitution{}), Kata{
		ID:          "ZC1012",
		Title:       "Use `$(command)` instead of backticks for command substitution",
		Description: "The `$(command)` syntax is generally preferred over backticks `` `command` `` for command substitution. It's easier to read, can be nested, and avoids issues with backslashes.",
		Check:       checkZC1012,
	})
}

func checkZC1012(node ast.Node) []Violation {
	violations := []Violation{}

	if cmdSub, ok := node.(*ast.CommandSubstitution); ok {
		violations = append(violations, Violation{
			KataID:  "ZC1012",
			Message: "Use `$(command)` instead of backticks for command substitution.",
			Line:    cmdSub.Token.Line,
			Column:  cmdSub.Token.Column,
		})
	}

	return violations
}
