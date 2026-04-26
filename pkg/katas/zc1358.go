// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1358",
		Title:    "Use `${PWD:P}` instead of `pwd -P` for physical current directory",
		Severity: SeverityStyle,
		Description: "`pwd -P` resolves symlinks to the physical path. Zsh's `${PWD:P}` modifier " +
			"does the same without spawning the external — the `P` modifier returns the " +
			"canonical (absolute, symlink-resolved) form.",
		Check: checkZC1358,
	})
}

func checkZC1358(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "pwd" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-P" {
			return []Violation{{
				KataID: "ZC1358",
				Message: "Use `${PWD:P}` instead of `pwd -P` — the `P` modifier resolves symlinks " +
					"and returns the canonical path without spawning an external.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}
