// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1315",
		Title:    "Avoid `$BASH_COMPAT` — use `emulate` for compatibility in Zsh",
		Severity: SeverityWarning,
		Description: "`$BASH_COMPAT` sets Bash compatibility level. Zsh uses `emulate` " +
			"to control compatibility mode (e.g., `emulate -L sh` for POSIX mode).",
		Check: checkZC1315,
	})
}

func checkZC1315(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident == nil {
		return nil
	}

	if ident.Value != "$BASH_COMPAT" && ident.Value != "BASH_COMPAT" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1315",
		Message: "Avoid `$BASH_COMPAT` in Zsh — use `emulate` for shell compatibility mode instead.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityWarning,
	}}
}
