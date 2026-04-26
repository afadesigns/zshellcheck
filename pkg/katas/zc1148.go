// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1148",
		Title:    "Use `compdef` instead of `compctl` for completions",
		Severity: SeverityInfo,
		Description: "`compctl` is the old Zsh completion system. " +
			"Use `compdef` with the new completion system (`compsys`) for modern Zsh.",
		Check: checkZC1148,
	})
}

func checkZC1148(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "compctl" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1148",
		Message: "Use `compdef` instead of `compctl`. The `compctl` system is deprecated; " +
			"use `compinit` and `compdef` for modern Zsh completions.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityInfo,
	}}
}
