// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1406",
		Title:    "Prefer Zsh `zargs -P N` autoload over `xargs -P N` for parallel execution",
		Severity: SeverityStyle,
		Description: "Zsh provides `zargs` (loaded via `autoload -Uz zargs`) — a native equivalent " +
			"of `xargs` with parallel execution via `-P`. It keeps variables and functions in " +
			"scope (unlike xargs) and avoids the utility-quoting surprises of `xargs`.",
		Check: checkZC1406,
	})
}

func checkZC1406(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "xargs" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-P" || v == "--max-procs" ||
			(len(v) > 2 && v[:2] == "-P") {
			_ = i
			return []Violation{{
				KataID: "ZC1406",
				Message: "Consider `zargs -P N` (autoload -Uz zargs) instead of `xargs -P N`. " +
					"Parallel execution with Zsh functions in scope — no subshell-per-item.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}
