package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1333",
		Title:    "Avoid `$TIMEFORMAT` — use `$TIMEFMT` in Zsh",
		Severity: SeverityInfo,
		Description: "`$TIMEFORMAT` is the Bash variable for customizing `time` output. " +
			"Zsh uses `$TIMEFMT` for the same purpose, with different format specifiers.",
		Check: checkZC1333,
	})
}

func checkZC1333(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value != "$TIMEFORMAT" && ident.Value != "TIMEFORMAT" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1333",
		Message: "Avoid `$TIMEFORMAT` in Zsh — use `$TIMEFMT` instead. Format specifiers differ between Bash and Zsh.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityInfo,
	}}
}
