// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1414",
		Title:    "Beware `hash -d` — Bash deletes from hash table, Zsh defines named directory",
		Severity: SeverityError,
		Description: "The `-d` flag has opposite meanings across shells: Bash `hash -d NAME` " +
			"removes `NAME` from the command-hash table. Zsh `hash -d NAME=PATH` **defines** a " +
			"named directory (`~NAME` expansion). A Bash script ported to Zsh breaks silently " +
			"when `hash -d ls` is interpreted as defining `~ls`.",
		Check: checkZC1414,
	})
}

func checkZC1414(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "hash" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-d" {
			return []Violation{{
				KataID: "ZC1414",
				Message: "`hash -d` has opposite semantics in Bash (delete) vs Zsh (define " +
					"named directory). Use `unhash cmd` for Zsh command-hash removal, or " +
					"`hash -d NAME=/path` for named-directory definition.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}

	return nil
}
