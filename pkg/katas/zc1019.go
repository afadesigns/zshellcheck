package katas

import (
	"reflect"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(reflect.TypeOf(&ast.SimpleCommand{}), Kata{
		ID:          "ZC1019",
		Title:       "Use `whence` instead of `which`",
		Description: "The `which` command is an external command and may not be available on all systems. The `whence` command is a built-in Zsh command that provides a more reliable and consistent way to find the location of a command.",
		Check:       checkZC1019,
	})
}

func checkZC1019(node ast.Node) []Violation {
	violations := []Violation{}

	if cmd, ok := node.(*ast.SimpleCommand); ok {
		if ident, ok := cmd.Name.(*ast.Identifier); ok {
			if ident.Value == "which" {
				violations = append(violations, Violation{
					KataID:  "ZC1019",
					Message: "Use `whence` instead of `which`.",
					Line:    ident.Token.Line,
					Column:  ident.Token.Column,
				})
			}
		}
	}

	return violations
}
