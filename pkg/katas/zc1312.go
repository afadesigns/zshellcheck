// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1312",
		Title:    "Avoid `compgen` command — use `compadd` in Zsh",
		Severity: SeverityWarning,
		Description: "`compgen` is a Bash builtin for generating completions. " +
			"Zsh uses `compadd` and the completion system functions for adding " +
			"completion candidates.",
		Check: checkZC1312,
	})
}

func checkZC1312(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "compgen" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1312",
		Message: "Avoid `compgen` in Zsh — it is a Bash builtin. Use `compadd` or Zsh completion functions instead.",
		Line:    cmd.Token.Line,
		Column:  cmd.Token.Column,
		Level:   SeverityWarning,
	}}
}
