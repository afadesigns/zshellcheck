// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1877",
		Title:    "Warn on `unsetopt SHORT_LOOPS` — short-form `for`/`while` bodies stop parsing",
		Severity: SeverityWarning,
		Description: "`SHORT_LOOPS` is on in Zsh by default: the compact forms `for x in *.log; " +
			"print $x`, `while true; print .`, and `repeat 3 sleep 1` parse with an implicit " +
			"single-command body. Turning the option off reverts to POSIX-shell parsing, " +
			"which demands an explicit `do ... done` or `{ ... }` block. Every subsequent " +
			"short-form loop raises a parse error (`parse error near '\\n'`), and the " +
			"behaviour is global so even helper files sourced later fall over. Keep the " +
			"option on; if you genuinely need POSIX-strict parsing, scope inside a function " +
			"with `setopt LOCAL_OPTIONS; unsetopt SHORT_LOOPS`.",
		Check: checkZC1877,
	})
}

func checkZC1877(node ast.Node) []Violation {
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
			if zc1877IsShortLoops(arg.String()) {
				return zc1877Hit(cmd, "unsetopt "+arg.String())
			}
		}
	case "setopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
			if norm == "NOSHORTLOOPS" {
				return zc1877Hit(cmd, "setopt "+v)
			}
		}
	}
	return nil
}

func zc1877IsShortLoops(v string) bool {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	return norm == "SHORTLOOPS"
}

func zc1877Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1877",
		Message: "`" + where + "` disables short-form loops — `for f in *.log; " +
			"print $f` raises a parse error. Keep the option on; scope inside " +
			"a function with `LOCAL_OPTIONS` if POSIX-strict parsing is " +
			"really needed.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
