// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1177",
		Title:    "Avoid `id -u` — use Zsh `$UID` or `$EUID`",
		Severity: SeverityStyle,
		Description: "Zsh provides `$UID` and `$EUID` as built-in variables for user/effective " +
			"user ID. Avoid spawning `id` for simple UID checks.",
		Check: checkZC1177,
	})
}

func checkZC1177(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "id" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-u" || val == "-un" {
			return []Violation{{
				KataID: "ZC1177",
				Message: "Use `$UID` or `$EUID` instead of `id -u`. " +
					"Zsh provides these as built-in variables.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}
