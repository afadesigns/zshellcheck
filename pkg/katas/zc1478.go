// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1478",
		Title:    "Avoid `mktemp -u` — returns a name without creating the file (TOCTOU)",
		Severity: SeverityWarning,
		Description: "`mktemp -u` allocates a unique name but does not create the file, leaving " +
			"a classic time-of-check to time-of-use race: a second process (possibly attacker- " +
			"controlled on a multi-user host or shared CI runner) can claim the name before you " +
			"redirect into it. Drop `-u` and operate on the file `mktemp` creates for you, or " +
			"use `mktemp -d` if you need a directory path.",
		Check: checkZC1478,
	})
}

func checkZC1478(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "mktemp" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-u" || v == "--dry-run" {
			return []Violation{{
				KataID: "ZC1478",
				Message: "`mktemp -u` returns a unique name but does not create the file — " +
					"TOCTOU race. Let `mktemp` create the file (or use `-d` for a directory).",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
