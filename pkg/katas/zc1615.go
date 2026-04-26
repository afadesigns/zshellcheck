// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1615",
		Title:    "Style: use Zsh `$EPOCHREALTIME` / `$epochtime` instead of `date \"+%s.%N\"`",
		Severity: SeverityStyle,
		Description: "Zsh's `zsh/datetime` module exposes `$EPOCHREALTIME` (scalar with " +
			"fractional seconds) and `$epochtime` (two-element array of seconds and " +
			"nanoseconds). Both read straight from `clock_gettime(CLOCK_REALTIME)` without " +
			"forking `date`. On a hot path the builtin is dramatically faster and avoids " +
			"subshell process-startup overhead. Autoload the module once with `zmodload " +
			"zsh/datetime`.",
		Check: checkZC1615,
	})
}

func checkZC1615(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "date" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if !strings.Contains(v, "%s") {
			continue
		}
		if strings.Contains(v, "%N") {
			return []Violation{{
				KataID: "ZC1615",
				Message: "`date " + v + "` forks for sub-second time. Use Zsh " +
					"`$EPOCHREALTIME` / `$epochtime` from `zmodload zsh/datetime`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}
	return nil
}
