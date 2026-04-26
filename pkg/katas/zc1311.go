// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1311",
		Title:    "Avoid `complete` command — use `compdef` in Zsh",
		Severity: SeverityWarning,
		Description: "`complete` is a Bash builtin for registering tab completions. " +
			"Zsh uses `compdef` for completion registration and the `compctl` " +
			"legacy interface. Use `compdef` for the modern Zsh completion system.",
		Check: checkZC1311,
	})
}

func checkZC1311(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "complete" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1311",
		Message: "Avoid `complete` in Zsh — it is a Bash builtin. Use `compdef` for Zsh completion registration.",
		Line:    cmd.Token.Line,
		Column:  cmd.Token.Column,
		Level:   SeverityWarning,
	}}
}
