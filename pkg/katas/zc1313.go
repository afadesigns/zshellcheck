package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1313",
		Title:    "Avoid `$BASH_ALIASES` — use Zsh `aliases` hash",
		Severity: SeverityWarning,
		Description: "`$BASH_ALIASES` is a Bash associative array of defined aliases. " +
			"Zsh provides the `aliases` associative array for the same purpose.",
		Check: checkZC1313,
	})
}

func checkZC1313(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value != "$BASH_ALIASES" && ident.Value != "BASH_ALIASES" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1313",
		Message: "Avoid `$BASH_ALIASES` in Zsh — use the `aliases` associative array instead.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityWarning,
	}}
}
