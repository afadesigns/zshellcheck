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

var (
	zc1481UnsetVars  = map[string]struct{}{"HISTFILE": {}, "HISTSIZE": {}, "SAVEHIST": {}, "HISTCMD": {}}
	zc1481EmptyHist  = map[string]struct{}{"": {}, "/dev/null": {}, "''": {}, `""`: {}}
	zc1481ZeroAssign = map[string]struct{}{"HISTSIZE=0": {}, "SAVEHIST=0": {}}
)

func checkZC1481(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	switch CommandIdentifier(cmd) {
	case "unset":
		if hit := zc1481UnsetHit(cmd); hit != "" {
			return zc1481Violation(cmd, "unset "+hit)
		}
	case "export", "typeset":
		if hit := zc1481AssignHit(cmd); hit != "" {
			return zc1481Violation(cmd, hit)
		}
	}
	return nil
}

func zc1481UnsetHit(cmd *ast.SimpleCommand) string {
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if _, hit := zc1481UnsetVars[v]; hit {
			return v
		}
	}
	return ""
}

func zc1481AssignHit(cmd *ast.SimpleCommand) string {
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if strings.HasPrefix(v, "HISTFILE=") {
			val := strings.TrimPrefix(v, "HISTFILE=")
			if _, hit := zc1481EmptyHist[val]; hit {
				return v
			}
		}
		if _, hit := zc1481ZeroAssign[v]; hit {
			return v
		}
	}
	return ""
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
