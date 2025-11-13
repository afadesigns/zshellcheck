package katas

import (
	"reflect"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(reflect.TypeOf(&ast.CallExpression{}), Kata{
		ID:          "ZC1005",
		Title:       "Use whence instead of which",
		Description: "The `which` command is an external command and may not be available on all systems. The `whence` command is a built-in Zsh command that provides a more reliable and consistent way to find the location of a command.",
		Check:       checkZC1005,
	})
}

func checkZC1005(node ast.Node) []Violation {
	violations := []Violation{}

	if callExpr, ok := node.(*ast.CallExpression); ok {
		if ident, ok := callExpr.Function.(*ast.Identifier); ok {
			if ident.Value == "which" {
				violations = append(violations, Violation{
					KataID:  "ZC1005",
					Message: "Use whence instead of which. The `whence` command is a built-in Zsh command that provides a more reliable and consistent way to find the location of a command.",
					Line:    callExpr.Token.Line,
					Column:  callExpr.Token.Column,
				})
			}
		}
	}

	return violations
}
