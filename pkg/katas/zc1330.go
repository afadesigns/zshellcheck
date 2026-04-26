// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1330",
		Title:    "Avoid `$INPUTRC` — use `bindkey` in Zsh",
		Severity: SeverityInfo,
		Description: "`$INPUTRC` points to the readline configuration file in Bash. " +
			"Zsh uses `bindkey` and ZLE widgets for key binding configuration, " +
			"not readline.",
		Check: checkZC1330,
	})
}

func checkZC1330(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident == nil {
		return nil
	}

	if ident.Value != "$INPUTRC" && ident.Value != "INPUTRC" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1330",
		Message: "Avoid `$INPUTRC` in Zsh — Zsh uses `bindkey` and ZLE, not readline. `INPUTRC` is Bash-specific.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityInfo,
	}}
}
