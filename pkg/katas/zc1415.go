// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1415",
		Title:    "Prefer Zsh `TRAPZERR` function over `trap 'cmd' ERR`",
		Severity: SeverityInfo,
		Description: "Both Bash and Zsh accept `trap 'cmd' ERR`, but Zsh's idiomatic form is the " +
			"named function `TRAPZERR`: `TRAPZERR() { echo \"err at $LINENO\"; }`. The named " +
			"function receives `$1` = signal and is easier to compose than an inline string.",
		Check: checkZC1415,
	})
}

func checkZC1415(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "trap" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "ERR" || v == "ZERR" {
			return []Violation{{
				KataID: "ZC1415",
				Message: "Prefer Zsh `TRAPZERR() { ... }` function over `trap 'cmd' ERR`. " +
					"The named-function form is more idiomatic and composable in Zsh.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityInfo,
			}}
		}
	}

	return nil
}
