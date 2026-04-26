// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1309",
		Title:    "Avoid `$BASH_COMMAND` — not available in Zsh",
		Severity: SeverityWarning,
		Description: "`$BASH_COMMAND` contains the currently executing command in Bash. " +
			"Zsh does not provide a direct equivalent. Use `$ZSH_DEBUG_CMD` in " +
			"debug traps or restructure the logic.",
		Check: checkZC1309,
	})
}

func checkZC1309(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident == nil {
		return nil
	}

	if ident.Value != "$BASH_COMMAND" && ident.Value != "BASH_COMMAND" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1309",
		Message: "Avoid `$BASH_COMMAND` in Zsh — it is undefined. Use `$ZSH_DEBUG_CMD` in debug traps if needed.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityWarning,
	}}
}
