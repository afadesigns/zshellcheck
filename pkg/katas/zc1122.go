package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:    "ZC1122",
		Title: "Use `$USER` instead of `whoami`",
		Description: "Zsh provides `$USER` as a built-in variable containing the current username. " +
			"Avoid spawning `whoami` as an external process.",
		Severity: SeverityStyle,
		Check:    checkZC1122,
	})
}

func checkZC1122(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok || ident == nil || ident.Value != "whoami" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1122",
		Message: "Use `$USER` instead of `whoami`. " +
			"Zsh maintains `$USER` as a built-in variable, avoiding an external process.",
		Line:   ident.Token.Line,
		Column: ident.Token.Column,
		Level:  SeverityStyle,
	}}
}
