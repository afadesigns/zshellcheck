// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1623",
		Title:    "Warn on `kill -STOP PID` / `pkill -STOP` — target halts until `kill -CONT` runs",
		Severity: SeverityWarning,
		Description: "Sending SIGSTOP halts the target process until SIGCONT arrives. If the " +
			"script fails, is killed, or exits before the resume, the target stays paused " +
			"indefinitely — consuming memory, holding locks, blocking its dependents. Wrap " +
			"every `kill -STOP $PID` with `trap \"kill -CONT $PID\" EXIT` (or an explicit " +
			"cleanup path) so the resume fires even on failure. Prefer `kill -TSTP` if the " +
			"target can handle it (the user-space tstop that the process can ignore).",
		Check: checkZC1623,
	})
}

func checkZC1623(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "kill" && ident.Value != "pkill" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-STOP" || v == "-SIGSTOP" || v == "-19" {
			return zc1623Hit(cmd)
		}
		if v == "-s" && i+1 < len(cmd.Arguments) {
			sig := cmd.Arguments[i+1].String()
			if sig == "STOP" || sig == "SIGSTOP" || sig == "19" {
				return zc1623Hit(cmd)
			}
		}
	}
	return nil
}

func zc1623Hit(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1623",
		Message: "`kill -STOP` halts the target until SIGCONT arrives. Pair every STOP " +
			"with `trap \"kill -CONT PID\" EXIT` so the resume fires even on failure.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
