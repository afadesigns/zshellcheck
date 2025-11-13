package katas

import (
	"reflect"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(reflect.TypeOf(&ast.SimpleCommand{}), Kata{
		ID:          "ZC1030",
		Title:       "Use `printf` instead of `echo`",
		Description: "The `echo` command's behavior can be inconsistent across different shells and environments, especially with flags and escape sequences. `printf` provides more reliable and portable string formatting.",
		Check:       checkZC1030,
	})
}

func checkZC1030(node ast.Node) []Violation {
	violations := []Violation{}

	if cmd, ok := node.(*ast.SimpleCommand); ok {
		if name, ok := cmd.Name.(*ast.Identifier); ok && name.Value == "echo" {
			violations = append(violations, Violation{
				KataID:  "ZC1030",
				Message: "Use `printf` for more reliable and portable string formatting instead of `echo`.",
				Line:    name.Token.Line,
				Column:  name.Token.Column,
			})
		}
	}

	return violations
}
