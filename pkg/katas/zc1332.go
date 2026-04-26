// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1332",
		Title:    "Avoid `$GLOBIGNORE` — use `setopt EXTENDED_GLOB` in Zsh",
		Severity: SeverityInfo,
		Description: "`$GLOBIGNORE` is a Bash variable for excluding patterns from glob expansion. " +
			"Zsh uses `setopt EXTENDED_GLOB` with the `~` (exclusion) operator or " +
			"`setopt NULL_GLOB` for different glob behavior.",
		Check: checkZC1332,
	})
}

func checkZC1332(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident == nil {
		return nil
	}

	if ident.Value != "$GLOBIGNORE" && ident.Value != "GLOBIGNORE" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1332",
		Message: "Avoid `$GLOBIGNORE` in Zsh — use `setopt EXTENDED_GLOB` with `~` operator for glob exclusion.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityInfo,
	}}
}
