package katas

import (
	"reflect"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(reflect.TypeOf(&ast.LetStatement{}), Kata{
		ID:          "ZC1029",
		Title:       "Use `$((...))` for arithmetic expansion",
		Description: "The `$((...))` syntax is the modern, recommended way to perform arithmetic expansion. It is more readable and can be nested easily, unlike `let`.",
		Check:       checkZC1029,
	})
}

func checkZC1029(node ast.Node) []Violation {
	violations := []Violation{}

	if let, ok := node.(*ast.LetStatement); ok {
		violations = append(violations, Violation{
			KataID:  "ZC1029",
			Message: "Use `$((...))` for arithmetic expansion instead of `let`.",
			Line:    let.Token.Line,
			Column:  let.Token.Column,
		})
	}

	return violations
}
