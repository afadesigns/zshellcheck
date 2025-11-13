package katas

import (
	"reflect"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(reflect.TypeOf(&ast.SimpleCommand{}), Kata{
		ID:          "ZC1007",
		Title:       "Avoid using `chmod 777`",
		Description: "Using `chmod 777` is a security risk as it gives read, write, and execute permissions to everyone. It's better to use more restrictive permissions.",
		Check:       checkZC1007,
	})
}

func checkZC1007(node ast.Node) []Violation {
	violations := []Violation{}

	if cmd, ok := node.(*ast.SimpleCommand); ok {
		if ident, ok := cmd.Name.(*ast.Identifier); ok {
			if ident.Value == "chmod" {
				for _, arg := range cmd.Arguments {
					switch v := arg.(type) {
					case *ast.Identifier:
						if v.Value == "777" {
							violations = append(violations, Violation{
								KataID:  "ZC1007",
								Message: "Avoid using `chmod 777`. It is a security risk.",
								Line:    ident.Token.Line,
								Column:  ident.Token.Column,
							})
						}
					case *ast.IntegerLiteral:
						if v.Token.Literal == "777" {
							violations = append(violations, Violation{
								KataID:  "ZC1007",
								Message: "Avoid using `chmod 777`. It is a security risk.",
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
