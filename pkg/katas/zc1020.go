package katas

import (
	"reflect"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(reflect.TypeOf(&ast.SimpleCommand{}), Kata{
		ID:          "ZC1020",
		Title:       "Use `[[ ... ]]` for tests instead of `test`",
		Description: "The `test` command is an external command and may not be available on all systems. The `[[...]]` construct is a Zsh keyword, offering safer and more powerful conditional expressions than the traditional `test` command.",
		Check:       checkZC1020,
	})
}

func checkZC1020(node ast.Node) []Violation {
	violations := []Violation{}

	if cmd, ok := node.(*ast.SimpleCommand); ok {
		if ident, ok := cmd.Name.(*ast.Identifier); ok {
			if ident.Value == "test" {
				violations = append(violations, Violation{
					KataID:  "ZC1020",
					Message: "Use `[[ ... ]]` for tests instead of `test`.",
					Line:    ident.Token.Line,
					Column:  ident.Token.Column,
				})
			}
		}
	}

	return violations
}
