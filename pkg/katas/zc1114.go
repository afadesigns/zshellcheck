// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1114",
		Title: "Consider Zsh `=(...)` for temporary files",
		Description: "Zsh `=(cmd)` creates a temporary file with the command output that is automatically " +
			"cleaned up. Consider this instead of manual `mktemp` and cleanup patterns.",
		Severity: SeverityStyle,
		Check:    checkZC1114,
	})
}

func checkZC1114(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "mktemp" {
		return nil
	}

	// Skip mktemp -d (directory creation — no Zsh equivalent)
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-d" {
			return nil
		}
	}

	return []Violation{{
		KataID: "ZC1114",
		Message: "Consider using Zsh `=(cmd)` for temporary files instead of `mktemp`. " +
			"Zsh auto-cleans temporary files created with `=(...)` process substitution.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}
