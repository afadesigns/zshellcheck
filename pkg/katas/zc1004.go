package katas

import (
	"reflect"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(reflect.TypeOf(&ast.CallExpression{}), Kata{
		ID:          "ZC1004",
		Title:       "Use print instead of echo -e for escaped strings",
		Description: "The `echo -e` command is not portable and can have unexpected behavior in Zsh. The `print` command is a built-in Zsh command that provides a more reliable and consistent way to print escaped strings.",
		Check:       checkZC1004,
	})
}

func checkZC1004(node ast.Node) []Violation {
	violations := []Violation{}

	if callExpr, ok := node.(*ast.CallExpression); ok {
		if ident, ok := callExpr.Function.(*ast.Identifier); ok {
			if ident.Value == "echo" {
				for _, arg := range callExpr.Arguments {
					if str, ok := arg.(*ast.StringLiteral); ok {
						if str.Value == "-e" {
							violations = append(violations, Violation{
								KataID:  "ZC1004",
								Message: "Use print instead of echo -e for escaped strings. The `print` command is a built-in Zsh command that provides a more reliable and consistent way to print escaped strings.",
								Line:    callExpr.Token.Line,
								Column:  callExpr.Token.Column,
							})
						}
					}
				}
			}
		}
	}

	return violations
}
