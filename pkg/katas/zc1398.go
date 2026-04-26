// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1398",
		Title:    "Avoid `$PROMPT_DIRTRIM` — use Zsh `%N~` prompt modifier",
		Severity: SeverityWarning,
		Description: "Bash's `$PROMPT_DIRTRIM` limits the number of directory components shown " +
			"in `\\w`. Zsh has no such variable; use the `%N~` prompt escape (N is component " +
			"count) or `%/` / `%~` with precmd adjustments for Zsh-native directory truncation.",
		Check: checkZC1398,
	})
}

func checkZC1398(node ast.Node) []Violation {
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
		if strings.Contains(v, "PROMPT_DIRTRIM") {
			return []Violation{{
				KataID: "ZC1398",
				Message: "`$PROMPT_DIRTRIM` is Bash-only. Use the Zsh prompt escape `%N~` " +
					"(N = number of path components to keep) for directory truncation.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}

	return nil
}
