package katas

import (
	"reflect"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(reflect.TypeOf(&ast.BracketExpression{}), Kata{
		ID:    "ZC1010",
		Title: "Use `[[ ... ]]` instead of `[ ... ]`",
		Description: "The `[[ ... ]]` construct is a Zsh keyword and is generally safer and more powerful " +
			"than the `[ ... ]` command (which is an alias for `test`). `[[ ... ]]` prevents word " +
			"splitting and pathname expansion, and supports advanced features like regex matching.",
		Check: checkZC1010,
	})
}

func checkZC1010(node ast.Node) []Violation {
	violations := []Violation{}

	if bracketExpr, ok := node.(*ast.BracketExpression); ok {
		violations = append(violations, Violation{
			KataID:  "ZC1010",
			Message: "Use `[[ ... ]]` instead of `[ ... ]` for safer and more powerful tests.",
			Line:    bracketExpr.Token.Line,
			Column:  bracketExpr.Token.Column,
		})
	}

	return violations
}
