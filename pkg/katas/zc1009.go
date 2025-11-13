package katas

import (
	"reflect"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(reflect.TypeOf(&ast.SimpleCommand{}), Kata{
		ID:          "ZC1009",
		Title:       "Use `((...))` for C-style arithmetic",
		Description: "The `((...))` construct in Zsh allows for C-style arithmetic. It is generally more efficient and readable than using `expr` or other external commands for arithmetic.",
		Check:       checkZC1009,
	})
}

func checkZC1009(node ast.Node) []Violation {
	violations := []Violation{}

	if cmd, ok := node.(*ast.SimpleCommand); ok {
		if ident, ok := cmd.Name.(*ast.Identifier); ok {
			if ident.Value == "expr" {
				violations = append(violations, Violation{
					KataID:  "ZC1009",
					Message: "Use `((...))` for C-style arithmetic instead of `expr`.",
					Line:    ident.Token.Line,
					Column:  ident.Token.Column,
				})
			}
		}
	}

	return violations
}
