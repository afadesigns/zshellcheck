// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1143",
		Title:    "Avoid `set -e` — use explicit error handling",
		Severity: SeverityInfo,
		Description: "`set -e` (errexit) has surprising behavior in Zsh with conditionals, " +
			"pipes, and subshells. Use explicit `|| return` or `|| exit` for reliable error handling.",
		Check: checkZC1143,
	})
}

func checkZC1143(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "set" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-e" || val == "-o" {
			// Check for set -o errexit pattern
			if val == "-e" {
				return []Violation{{
					KataID: "ZC1143",
					Message: "Avoid `set -e`. It has surprising behavior with conditionals and subshells in Zsh. " +
						"Use explicit error handling with `cmd || return 1` instead.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityInfo,
				}}
			}
		}
	}

	return nil
}
