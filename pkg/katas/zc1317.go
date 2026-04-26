// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1317",
		Title:    "Avoid `$BASH_ENV` — use `$ZDOTDIR` and `$ENV` in Zsh",
		Severity: SeverityInfo,
		Description: "`$BASH_ENV` specifies a startup file for non-interactive Bash shells. " +
			"Zsh uses `$ZDOTDIR` to locate `.zshrc` and related files, and `$ENV` " +
			"for POSIX-compatible startup.",
		Check: checkZC1317,
	})
}

func checkZC1317(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident == nil {
		return nil
	}

	if ident.Value != "$BASH_ENV" && ident.Value != "BASH_ENV" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1317",
		Message: "Avoid `$BASH_ENV` in Zsh — use `$ZDOTDIR` for Zsh startup file locations instead.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityInfo,
	}}
}
