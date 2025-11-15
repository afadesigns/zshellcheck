package katas

import (
	"reflect"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(reflect.TypeOf(&ast.ExpressionStatement{}), Kata{
		ID:          "ZC1034",
		Title:       "Use `command -v` instead of `which`",
		Description: "`which` is an external command and may not be available or consistent across all " +
			"systems. `command -v` is a POSIX standard and a shell builtin, making it more portable " +
			"and reliable for checking if a command exists.",
		Check:       checkZC1034,
	})
}

func checkZC1034(node ast.Node) []Violation {
	violations := []Violation{}

	if es, ok := node.(*ast.ExpressionStatement); ok {
		if cmd, ok := es.Expression.(*ast.SimpleCommand); ok {
			if cmd.Name.TokenLiteral() == "which" {
				violations = append(violations, Violation{
					KataID:  "ZC1034",
					Message: "Use `command -v` instead of `which` for portability.",
					Line:    cmd.Token.Line,
					Column:  cmd.Token.Column,
				})
			}
		}
	}

	return violations
}
