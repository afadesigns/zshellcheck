package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1328",
		Title:    "Avoid `$HISTCONTROL` — use Zsh `setopt` history options",
		Severity: SeverityInfo,
		Description: "`$HISTCONTROL` is a Bash variable controlling history deduplication. " +
			"Zsh uses `setopt HIST_IGNORE_DUPS`, `HIST_IGNORE_ALL_DUPS`, and " +
			"`HIST_IGNORE_SPACE` for the same functionality.",
		Check: checkZC1328,
	})
}

func checkZC1328(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident == nil {
		return nil
	}

	if ident.Value != "$HISTCONTROL" && ident.Value != "HISTCONTROL" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1328",
		Message: "Avoid `$HISTCONTROL` in Zsh — use `setopt HIST_IGNORE_DUPS` and related options instead.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityInfo,
	}}
}
