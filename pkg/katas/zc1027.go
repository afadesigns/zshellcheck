package katas

import (
	"reflect"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(reflect.TypeOf(&ast.LetStatement{}), Kata{
		ID:          "ZC1027",
		Title:       "Use `$((...))` for arithmetic expansion",
		Description: "The `$((...))` syntax is the modern, recommended way to perform arithmetic expansion. " +
			"It is more readable and can be nested easily, unlike `let`.",
		Check:       checkZC1027,
	})
}

func checkZC1027(node ast.Node) []Violation {
	violations := []Violation{}

	if let, ok := node.(*ast.LetStatement); ok {
		violations = append(violations, Violation{
			KataID:  "ZC1027",
			Message: "Use `$((...))` for arithmetic expansion instead of `let`.",
			Line:    let.Token.Line,
			Column:  let.Token.Column,
		})
	}

	return violations
}
