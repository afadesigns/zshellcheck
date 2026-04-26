// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1359",
		Title:    "Avoid `id -Gn` — use Zsh `$groups` associative array",
		Severity: SeverityStyle,
		Description: "Zsh's `zsh/parameter` module exposes the `$groups` associative array mapping " +
			"group names to GIDs for the current process. Load with `zmodload zsh/parameter` " +
			"(often auto-loaded) and inspect `${(k)groups}` for names, avoiding an external " +
			"`id -Gn`/`groups` call.",
		Check: checkZC1359,
	})
}

func checkZC1359(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "id" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-Gn" || v == "-G" || v == "-gn" || v == "-g" {
			return []Violation{{
				KataID: "ZC1359",
				Message: "Avoid `id -Gn`/`-G`/`-gn`/`-g` — use Zsh `$groups` (names→gids assoc array) " +
					"or `$GID` for the primary group after `zmodload zsh/parameter`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}
