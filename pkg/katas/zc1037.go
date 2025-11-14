package katas

import (
	"reflect"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(reflect.TypeOf(&ast.PrefixExpression{}), Kata{
		ID:          "ZC1037",
		Title:       "Quote variable expansions",
		Description: "Unquoted variable expansions can lead to unexpected behavior due to word splitting and globbing. It is recommended to quote all variable expansions.",
		Check:       checkZC1037,
	})
}

func checkZC1037(node ast.Node) []Violation {
	violations := []Violation{}

	if prefix, ok := node.(*ast.PrefixExpression); ok {
		if prefix.Operator == "$" {
			if _, ok := prefix.Right.(*ast.Identifier); ok {
				violations = append(violations, Violation{
					KataID:  "ZC1037",
					Message: "Unquoted variable expansion. Quote to prevent word splitting and globbing.",
					Line:    prefix.Token.Line,
					Column:  prefix.Token.Column,
				})
			}
		}
	}

	return violations
}
