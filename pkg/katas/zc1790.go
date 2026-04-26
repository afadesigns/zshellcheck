// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1790",
		Title:    "Warn on `unsetopt PIPE_FAIL` — pipeline exit status reverts to last-command-only",
		Severity: SeverityWarning,
		Description: "With `PIPE_FAIL` off (the shell default), `cmd1 | cmd2 | cmd3` exits with " +
			"`cmd3`'s status; failures in `cmd1` and `cmd2` are silently dropped. " +
			"`unsetopt PIPE_FAIL` (or the equivalent `setopt NOPIPEFAIL`) mid-script turns a " +
			"previously-enabled error check back off — typically because a known-flaky pipe " +
			"stage was tripping `set -e`, and the author reached for the global off-switch. " +
			"Undo the change in a subshell (`( unsetopt pipefail; …; )`) or a function with " +
			"`emulate -L zsh; unsetopt pipefail` so the rest of the script keeps strict-pipe " +
			"error propagation.",
		Check: checkZC1790,
	})
}

func checkZC1790(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "unsetopt":
		for _, arg := range cmd.Arguments {
			if zc1790IsPipeFail(arg.String()) {
				return zc1790Hit(cmd, "unsetopt "+arg.String())
			}
		}
	case "setopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
			if norm == "NOPIPEFAIL" {
				return zc1790Hit(cmd, "setopt "+v)
			}
		}
	}
	return nil
}

func zc1790IsPipeFail(v string) bool {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	return norm == "PIPEFAIL"
}

func zc1790Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1790",
		Message: "`" + where + "` returns the shell to last-command-only pipeline exit — " +
			"`cmd1 | cmd2` now ignores `cmd1` failures. Scope the change to a subshell " +
			"or function with `emulate -L zsh` instead of flipping it globally.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
