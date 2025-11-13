package katas

import (
	"reflect"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(reflect.TypeOf(&ast.LetStatement{}), Kata{
		ID:          "ZC1008",
		Title:       "Use `\\$(())` for arithmetic operations",
		Description: "The `let` command is a shell builtin, but the `\\$(())` syntax is more portable and generally preferred for arithmetic operations in Zsh. It's also more powerful as it can be used in more contexts.",
		Check:       checkZC1008,
	})
}

func checkZC1008(node ast.Node) []Violation {
	violations := []Violation{}

	if letStmt, ok := node.(*ast.LetStatement); ok {
		violations = append(violations, Violation{
			KataID:  "ZC1008",
			Message: "Use `\\$(())` for arithmetic operations instead of `let`.",
			Line:    letStmt.Token.Line,
			Column:  letStmt.Token.Column,
		})
	}

	return violations
}
