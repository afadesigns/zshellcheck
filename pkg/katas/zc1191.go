package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1191",
		Title:    "Avoid `clear` command — use ANSI escape sequences",
		Severity: SeverityStyle,
		Description: "`clear` spawns an external process for screen clearing. " +
			"Use `print -n '\\e[2J\\e[H'` for faster terminal clearing.",
		Check: checkZC1191,
	})
}

func checkZC1191(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok || ident == nil || ident.Value != "clear" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1191",
		Message: "Use `print -n '\\e[2J\\e[H'` instead of `clear`. " +
			"ANSI escape sequences avoid spawning an external process.",
		Line:   ident.Token.Line,
		Column: ident.Token.Column,
		Level:  SeverityStyle,
	}}
}
