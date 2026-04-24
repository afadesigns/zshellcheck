package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1314",
		Title:    "Avoid `$BASH_LOADABLES_PATH` — not available in Zsh",
		Severity: SeverityWarning,
		Description: "`$BASH_LOADABLES_PATH` is a Bash variable for loadable builtin search paths. " +
			"Zsh has no equivalent; use `zmodload` with full module names instead.",
		Check: checkZC1314,
	})
}

func checkZC1314(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident == nil {
		return nil
	}

	if ident.Value != "$BASH_LOADABLES_PATH" && ident.Value != "BASH_LOADABLES_PATH" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1314",
		Message: "Avoid `$BASH_LOADABLES_PATH` in Zsh — it is undefined. Use `zmodload` with full module names.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityWarning,
	}}
}
