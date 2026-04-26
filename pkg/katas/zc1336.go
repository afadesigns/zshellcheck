// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1336",
		Title:    "Avoid `printenv` — use `typeset -x` or `export` in Zsh",
		Severity: SeverityStyle,
		Description: "`printenv` is an external command for listing environment variables. " +
			"Zsh provides `typeset -x` to list exported variables and `export` " +
			"to display them without spawning a subprocess.",
		Check: checkZC1336,
	})
}

func checkZC1336(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "printenv" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1336",
		Message: "Avoid `printenv` in Zsh — use `typeset -x` or `export` to list environment variables.",
		Line:    cmd.Token.Line,
		Column:  cmd.Token.Column,
		Level:   SeverityStyle,
	}}
}
