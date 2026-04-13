package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1321",
		Title:    "Avoid `$BASH_XTRACEFD` — not available in Zsh",
		Severity: SeverityWarning,
		Description: "`$BASH_XTRACEFD` redirects Bash xtrace output to a file descriptor. " +
			"Zsh does not have this variable. Use `exec 2>file` or redirect " +
			"stderr directly for trace output redirection.",
		Check: checkZC1321,
	})
}

func checkZC1321(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value != "$BASH_XTRACEFD" && ident.Value != "BASH_XTRACEFD" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1321",
		Message: "Avoid `$BASH_XTRACEFD` in Zsh — it is undefined. Redirect stderr directly for xtrace output.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityWarning,
	}}
}
