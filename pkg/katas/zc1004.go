package katas

import (
	"reflect"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(reflect.TypeOf(&ast.SimpleCommand{}), Kata{
		ID:          "ZC1004",
		Title:       "Use print instead of echo -e for escaped strings",
		Description: "The `echo -e` command is not portable and can have unexpected behavior in Zsh. The `print` command is a built-in Zsh command that provides a more reliable and consistent way to print escaped strings.",
		Check:       checkZC1004,
	})
}

func checkZC1004(node ast.Node) []Violation {
	violations := []Violation{}

	if cmd, ok := node.(*ast.SimpleCommand); ok {
		if ident, ok := cmd.Name.(*ast.Identifier); ok {
			if ident.Value == "echo" {
				for _, arg := range cmd.Arguments {
					if str, ok := arg.(*ast.StringLiteral); ok {
						if str.Value == "-e" {
							violations = append(violations, Violation{
								KataID:  "ZC1004",
								Message: "Use print instead of echo -e for escaped strings. The `print` command is a built-in Zsh command that provides a more reliable and consistent way to print escaped strings.",
								Line:    ident.Token.Line,
								Column:  ident.Token.Column,
							})
						}
					}
				}
			}
		}
	}

	return violations
}
