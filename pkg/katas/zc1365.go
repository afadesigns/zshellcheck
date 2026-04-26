// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1365",
		Title:    "Use Zsh `zstat` module instead of `stat -c` for file metadata",
		Severity: SeverityStyle,
		Description: "Zsh's `zsh/stat` module (loaded with `zmodload zsh/stat` — the command is " +
			"named `zstat`) exposes every `stat(2)` field natively: mtime, size, owner, group, " +
			"mode, links, etc. Avoid external `stat -c '%...'` invocations.",
		Check: checkZC1365,
	})
}

func checkZC1365(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "stat" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-c" || v == "--format" || v == "--printf" {
			return []Violation{{
				KataID: "ZC1365",
				Message: "Use Zsh `zmodload zsh/stat; zstat -H meta file` for file metadata instead " +
					"of `stat -c '%...'`. The associative array `meta` exposes every stat field.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}
