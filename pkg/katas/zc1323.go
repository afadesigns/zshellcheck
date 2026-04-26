// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1323",
		Title:    "Avoid `suspend` builtin — use `kill -STOP $$` in Zsh",
		Severity: SeverityWarning,
		Description: "`suspend` is a Bash builtin that suspends the shell. Zsh does not have " +
			"a `suspend` builtin. Use `kill -STOP $$` or Ctrl-Z for the same effect.",
		Check: checkZC1323,
	})
}

func checkZC1323(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "suspend" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1323",
		Message: "Avoid `suspend` in Zsh — it is a Bash builtin. Use `kill -STOP $$` if needed.",
		Line:    cmd.Token.Line,
		Column:  cmd.Token.Column,
		Level:   SeverityWarning,
	}}
}
