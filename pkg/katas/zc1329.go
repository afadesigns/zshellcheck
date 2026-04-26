// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1329",
		Title:    "Avoid `$HISTIGNORE` — use `zshaddhistory` hook in Zsh",
		Severity: SeverityInfo,
		Description: "`$HISTIGNORE` is a Bash variable for pattern-based history filtering. " +
			"Zsh uses the `zshaddhistory` hook function and `setopt HIST_IGNORE_SPACE` " +
			"for controlling which commands enter history.",
		Check: checkZC1329,
	})
}

func checkZC1329(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident == nil {
		return nil
	}

	if ident.Value != "$HISTIGNORE" && ident.Value != "HISTIGNORE" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1329",
		Message: "Avoid `$HISTIGNORE` in Zsh — use `zshaddhistory` hook for history filtering instead.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityInfo,
	}}
}
