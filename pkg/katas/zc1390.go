// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1390",
		Title:    "Avoid `$GROUPS[@]` — Zsh `$GROUPS` is a scalar, not an array",
		Severity: SeverityError,
		Description: "Bash's `$GROUPS` is an array of all group IDs the user belongs to, so " +
			"`${GROUPS[@]}` iterates them. In Zsh, `$GROUPS` is a scalar (primary GID). The " +
			"array of all group IDs is `$(groups)` output or `${(k)groups}` (if the " +
			"`zsh/parameter` module is loaded, `$groups` is an assoc array name→gid).",
		Check: checkZC1390,
	})
}

func checkZC1390(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "echo" && ident.Value != "print" && ident.Value != "printf" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if strings.Contains(v, "GROUPS[") || strings.Contains(v, "${GROUPS[") {
			return []Violation{{
				KataID: "ZC1390",
				Message: "Zsh `$GROUPS` is a scalar (primary GID), not an array. For all group " +
					"IDs use `${(k)groups}` (after `zmodload zsh/parameter`).",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}

	return nil
}
