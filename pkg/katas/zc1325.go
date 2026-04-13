package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1325",
		Title:    "Avoid `$PS0` — use `preexec` hook in Zsh",
		Severity: SeverityWarning,
		Description: "`$PS0` is a Bash 4.4+ prompt string displayed before command execution. " +
			"Zsh uses the `preexec` hook function for running code before each command.",
		Check: checkZC1325,
	})
}

func checkZC1325(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value != "$PS0" && ident.Value != "PS0" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1325",
		Message: "Avoid `$PS0` in Zsh — use the `preexec` hook function instead. `PS0` is Bash 4.4+ specific.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityWarning,
	}}
}
