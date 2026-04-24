package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1297",
		Title:    "Avoid `$BASH_SOURCE` — use `$0` or `${(%):-%x}` in Zsh",
		Severity: SeverityWarning,
		Description: "`$BASH_SOURCE` is a Bash-specific variable that does not exist in Zsh. " +
			"In Zsh, use `$0` inside a sourced file to get the script path, or " +
			"`${(%):-%x}` for the current file regardless of sourcing context.",
		Check: checkZC1297,
	})
}

func checkZC1297(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident == nil {
		return nil
	}

	if ident.Value != "$BASH_SOURCE" && ident.Value != "BASH_SOURCE" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1297",
		Message: "Avoid `$BASH_SOURCE` in Zsh — use `$0` or `${(%):-%x}` instead. `BASH_SOURCE` is Bash-specific.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityWarning,
	}}
}
