// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1299",
		Title:    "Avoid `$BASH_LINENO` — use `$funcfiletrace` in Zsh",
		Severity: SeverityWarning,
		Description: "`$BASH_LINENO` is a Bash-specific array that does not exist in Zsh. " +
			"Zsh provides `$funcfiletrace` as the equivalent, containing file:line " +
			"pairs for each call in the function stack.",
		Check: checkZC1299,
	})
}

func checkZC1299(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident == nil {
		return nil
	}

	if ident.Value != "$BASH_LINENO" && ident.Value != "BASH_LINENO" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1299",
		Message: "Avoid `$BASH_LINENO` in Zsh — use `$funcfiletrace` instead. `BASH_LINENO` is Bash-specific.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityWarning,
	}}
}
