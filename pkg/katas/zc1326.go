package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1326",
		Title:    "Avoid `$HISTTIMEFORMAT` — use `fc -li` in Zsh",
		Severity: SeverityInfo,
		Description: "`$HISTTIMEFORMAT` is a Bash variable for formatting history timestamps. " +
			"Zsh stores timestamps automatically when `EXTENDED_HISTORY` is set, " +
			"and displays them with `fc -li` or `history -i`.",
		Check: checkZC1326,
	})
}

func checkZC1326(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value != "$HISTTIMEFORMAT" && ident.Value != "HISTTIMEFORMAT" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1326",
		Message: "Avoid `$HISTTIMEFORMAT` in Zsh — use `setopt EXTENDED_HISTORY` and `fc -li` instead.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityInfo,
	}}
}
