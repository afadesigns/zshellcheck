// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1327",
		Title:    "Avoid `history -c` — Zsh uses different history management",
		Severity: SeverityWarning,
		Description: "`history -c` clears history in Bash. Zsh provides `fc -p` for pushing " +
			"history to a new file and `fc -P` for popping. Use `fc -W` to write and " +
			"`fc -R` to read history files.",
		Check: checkZC1327,
	})
}

func checkZC1327(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "history" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		// `-c` / `-d` are the destructive / anti-forensics flags and are
		// owned by ZC1487; this kata narrows to the Bash-only write/read
		// portability flags (`-w` / `-r` / `-a`).
		if val == "-w" || val == "-r" || val == "-a" {
			return []Violation{{
				KataID:  "ZC1327",
				Message: "Avoid `history " + val + "` in Zsh — Bash history flags differ. Use `fc` commands for Zsh history management.",
				Line:    cmd.Token.Line,
				Column:  cmd.Token.Column,
				Level:   SeverityWarning,
			}}
		}
	}

	return nil
}
