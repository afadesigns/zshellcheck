// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1530",
		Title:    "Warn on `pkill -f <pattern>` — matches full command line, easy to over-kill",
		Severity: SeverityWarning,
		Description: "`pkill -f` matches the pattern against the full command line, not just " +
			"the process name. A pattern like `-f server` also matches the `grep -- server` " +
			"in a user's shell history or any backup tool named `server-backup`. For routine " +
			"use, drop `-f` (matches process name only) or scope with `-U <uid>` / `-G " +
			"<gid>` / `-P <ppid>`. When you must match the command line, pin it with `^` / `$` " +
			"anchors in the pattern.",
		Check: checkZC1530,
	})
}

func checkZC1530(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "pkill" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-f" {
			return []Violation{{
				KataID: "ZC1530",
				Message: "`pkill -f` matches the full command line — easy to over-kill. Drop " +
					"`-f`, scope with `-U/-G/-P`, or anchor the pattern with ^/$.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
