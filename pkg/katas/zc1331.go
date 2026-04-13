package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1331",
		Title:    "Avoid `$BASH_REMATCH` — use `$match` array in Zsh",
		Severity: SeverityWarning,
		Description: "`$BASH_REMATCH` holds regex capture groups in Bash. Zsh stores " +
			"regex matches in the `$match` array (and `$MATCH` for the full match) " +
			"when using `=~` with `setopt BASH_REMATCH` disabled.",
		Check: checkZC1331,
	})
}

func checkZC1331(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value != "$BASH_REMATCH" && ident.Value != "BASH_REMATCH" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1331",
		Message: "Avoid `$BASH_REMATCH` in Zsh — use `$match` array and `$MATCH` for regex captures instead.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityWarning,
	}}
}
