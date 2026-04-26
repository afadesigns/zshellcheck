// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1417",
		Title:    "Prefer Zsh `TRAPRETURN` function over `trap 'cmd' RETURN`",
		Severity: SeverityInfo,
		Description: "Bash's `trap 'cmd' RETURN` runs `cmd` when a function returns. Zsh accepts " +
			"the `RETURN` signal name but the idiomatic form is a function named `TRAPRETURN`: " +
			"`TRAPRETURN() { print \"returning $?\"; }`.",
		Check: checkZC1417,
	})
}

func checkZC1417(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "trap" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "RETURN" {
			return []Violation{{
				KataID: "ZC1417",
				Message: "Prefer Zsh `TRAPRETURN() { ... }` function over `trap 'cmd' RETURN`. " +
					"Named-function form is more idiomatic in Zsh.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityInfo,
			}}
		}
	}

	return nil
}
