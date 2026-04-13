package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1300",
		Title:    "Avoid `$BASH_VERSINFO` — use `$ZSH_VERSION` in Zsh",
		Severity: SeverityWarning,
		Description: "`$BASH_VERSINFO` is a Bash-specific array containing version components. " +
			"In Zsh, use `$ZSH_VERSION` (string) or `${(s:.:)ZSH_VERSION}` to split " +
			"it into components for version comparison.",
		Check: checkZC1300,
	})
}

func checkZC1300(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value != "$BASH_VERSINFO" && ident.Value != "BASH_VERSINFO" &&
		ident.Value != "$BASH_VERSION" && ident.Value != "BASH_VERSION" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1300",
		Message: "Avoid Bash version variables in Zsh — use `$ZSH_VERSION` instead. Bash version variables are undefined in Zsh.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityWarning,
	}}
}
