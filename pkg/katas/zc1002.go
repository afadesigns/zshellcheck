package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.CommandSubstitutionNode, Kata{
		ID:          "ZC1002",
		Title:       "Use $(...) instead of backticks",
		Description: "Backticks are the old-style command substitution. $(...) is nesting-safe, easier to read, and generally preferred.",
		Check:       checkZC1002,
	})
}

func checkZC1002(node ast.Node) []Violation {
	violations := []Violation{}

	if cs, ok := node.(*ast.CommandSubstitution); ok {
		violations = append(violations, Violation{
			KataID:  "ZC1002",
			Message: "Use $(...) instead of backticks for command substitution. The `$(...)` syntax is more readable and can be nested easily.",
			Line:    cs.Token.Line,
			Column:  cs.Token.Column,
		})
	}

	return violations
}
