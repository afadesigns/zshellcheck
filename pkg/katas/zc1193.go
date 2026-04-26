// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1193",
		Title:    "Avoid `rm -i` in non-interactive scripts",
		Severity: SeverityWarning,
		Description: "`rm -i` prompts for confirmation which hangs in non-interactive scripts. " +
			"Remove the `-i` flag or use `rm -f` for scripts that run unattended.",
		Check: checkZC1193,
	})
}

func checkZC1193(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "rm" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-i" || val == "-ri" || val == "-ir" {
			return []Violation{{
				KataID: "ZC1193",
				Message: "Avoid `rm -i` in scripts — it prompts interactively and will hang " +
					"in non-interactive execution. Remove `-i` or use explicit checks instead.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}

	return nil
}
