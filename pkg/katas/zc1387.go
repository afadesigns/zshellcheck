// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1387",
		Title:    "Avoid `$SHELLOPTS` — Zsh uses `$options` associative array",
		Severity: SeverityWarning,
		Description: "Bash's `$SHELLOPTS` is a colon-separated list of set options. Zsh exposes " +
			"the same information via the `$options` associative array (keys are option names, " +
			"values are `on`/`off`). `$SHELLOPTS` is unset in Zsh.",
		Check: checkZC1387,
	})
}

func checkZC1387(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "echo" && ident.Value != "print" && ident.Value != "printf" && ident.Value != "export" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if strings.Contains(v, "SHELLOPTS") {
			return []Violation{{
				KataID: "ZC1387",
				Message: "`$SHELLOPTS` is Bash-only. In Zsh inspect `$options` (assoc array, " +
					"keys are option names) via `print -l ${(kv)options}`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}

	return nil
}
