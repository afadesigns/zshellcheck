// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1481",
		Title:    "Warn on `unset HISTFILE` / `export HISTFILE=/dev/null` — disables shell history",
		Severity: SeverityWarning,
		Description: "Disabling shell history (`unset HISTFILE`, `HISTFILE=/dev/null`, " +
			"`HISTSIZE=0`) is a classic stepping stone for hiding post-compromise activity. " +
			"Legitimate scripts almost never need this — if you are pasting a secret on the " +
			"command line, use `HISTCONTROL=ignorespace` and prefix the line with a space, or " +
			"read the value from a file / stdin.",
		Check: checkZC1481,
	})
}

func checkZC1481(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "unset":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			if v == "HISTFILE" || v == "HISTSIZE" || v == "SAVEHIST" || v == "HISTCMD" {
				return zc1481Violation(cmd, "unset "+v)
			}
		}
	case "export", "typeset":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			if strings.HasPrefix(v, "HISTFILE=") {
				val := strings.TrimPrefix(v, "HISTFILE=")
				if val == "" || val == "/dev/null" || val == "''" || val == `""` {
					return zc1481Violation(cmd, v)
				}
			}
			if v == "HISTSIZE=0" || v == "SAVEHIST=0" {
				return zc1481Violation(cmd, v)
			}
		}
	}
	return nil
}

func zc1481Violation(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1481",
		Message: "`" + what + "` disables shell history — textbook post-compromise tactic. " +
			"Legitimate alternative: `HISTCONTROL=ignorespace` plus leading-space prefix.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
