// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1887",
		Title:    "Warn on `setopt POSIX_TRAPS` — EXIT/ZERR traps change scope and no longer fire on function return",
		Severity: SeverityWarning,
		Description: "`POSIX_TRAPS` is off by default in Zsh. With it off, `trap cleanup EXIT` " +
			"inside a function fires when that function returns — the idiomatic Zsh way " +
			"to scope cleanup to a scope. Turning the option on reverts to POSIX-sh " +
			"semantics, where the EXIT trap only fires when the whole shell exits and " +
			"is shared across the entire process. Scripts that installed a cleanup trap " +
			"inside `do_work()` expecting it to run at each invocation now leak the " +
			"first trap's handler into everything after, and helpers that counted on " +
			"TRAPZERR / TRAPEXIT function-scoped behaviour silently skip. Keep the " +
			"option off at script level; if a specific line really needs POSIX-scope, " +
			"use `trap … EXIT` at top level and document it.",
		Check: checkZC1887,
	})
}

func checkZC1887(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "setopt":
		for _, arg := range cmd.Arguments {
			if zc1887IsPosixTraps(arg.String()) {
				return zc1887Hit(cmd, "setopt "+arg.String())
			}
		}
	case "unsetopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
			if norm == "NOPOSIXTRAPS" {
				return zc1887Hit(cmd, "unsetopt "+v)
			}
		}
	}
	return nil
}

func zc1887IsPosixTraps(v string) bool {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	return norm == "POSIXTRAPS"
}

func zc1887Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1887",
		Message: "`" + where + "` flips `trap ... EXIT` inside functions from " +
			"function-return to shell-exit scope — per-call cleanup leaks across " +
			"the whole shell, TRAPZERR helpers stop firing. Keep the option off.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
