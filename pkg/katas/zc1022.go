package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.LetStatementNode, Kata{
		ID:    "ZC1022",
		Title: "Use `$((...))` for arithmetic expansion",
		Description: "The `$((...))` syntax is the modern, recommended way to perform arithmetic expansion. " +
			"It is more readable and can be nested easily, unlike `let`.",
		Severity: SeverityStyle,
		Check:    checkZC1022,
		// Reuse ZC1013's `let NAME=EXPR` → `(( NAME = EXPR ))` rewrite.
		// For a standalone arithmetic statement the `(( ))` command
		// form is the right shape; the `$((...))` text in the message
		// reads as the broader "use Zsh arithmetic" recommendation.
		Fix: fixZC1013,
	})
}

func checkZC1022(node ast.Node) []Violation {
	violations := []Violation{}

	if let, ok := node.(*ast.LetStatement); ok {
		violations = append(violations, Violation{
			KataID:  "ZC1022",
			Message: "Use `$((...))` for arithmetic expansion instead of `let`.",
			Line:    let.Token.Line,
			Column:  let.Token.Column,
			Level:   SeverityStyle,
		})
	}

	return violations
}
