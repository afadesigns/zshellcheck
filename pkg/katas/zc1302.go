// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1302",
		Title:    "Avoid `help` builtin — use `run-help` or `man` in Zsh",
		Severity: SeverityInfo,
		Description: "The `help` command is a Bash builtin that displays builtin help. " +
			"Zsh does not have a `help` builtin. Use `run-help <command>` or " +
			"`man zshbuiltins` for Zsh builtin documentation.",
		Check: checkZC1302,
	})
}

func checkZC1302(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "help" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1302",
		Message: "Avoid `help` in Zsh — it is a Bash builtin. Use `run-help` or `man zshbuiltins` instead.",
		Line:    cmd.Token.Line,
		Column:  cmd.Token.Column,
		Level:   SeverityInfo,
	}}
}
