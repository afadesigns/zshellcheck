// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1716",
		Title:    "Use Zsh `$CPUTYPE` / `$MACHTYPE` instead of `uname -m` / `-p`",
		Severity: SeverityStyle,
		Description: "Zsh maintains `$CPUTYPE` (e.g. `x86_64`, `aarch64`) and `$MACHTYPE` " +
			"(the GNU triplet) as built-in parameters. Reading them is a constant-time " +
			"parameter expansion, while `uname -m` / `uname -p` forks an external for the " +
			"same answer. The Zsh values are populated at shell start from the same `uname(2)` " +
			"call, so they stay in lockstep with what `uname` would print.",
		Check: checkZC1716,
	})
}

func checkZC1716(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "uname" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-m" || v == "-p" {
			return []Violation{{
				KataID: "ZC1716",
				Message: "Use Zsh `$CPUTYPE` / `$MACHTYPE` instead of `uname " + v + "` — " +
					"parameter expansion avoids forking an external for an answer Zsh " +
					"already cached at startup.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}
	return nil
}
