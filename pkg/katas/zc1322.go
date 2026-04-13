package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1322",
		Title:    "Avoid `$COPROC` — Zsh coproc uses different syntax",
		Severity: SeverityWarning,
		Description: "`$COPROC` is a Bash array for coprocess file descriptors. " +
			"Zsh coprocesses use `coproc` keyword with different variable naming " +
			"and `read -p`/`print -p` for I/O.",
		Check: checkZC1322,
	})
}

func checkZC1322(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value != "$COPROC" && ident.Value != "COPROC" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1322",
		Message: "Avoid `$COPROC` in Zsh — Zsh coprocesses use `read -p`/`print -p` for I/O. `COPROC` is Bash-specific.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityWarning,
	}}
}
