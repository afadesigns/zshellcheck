// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1800",
		Title:    "Warn on `pg_ctl stop -m immediate` — abrupt shutdown skips checkpoint, forces WAL recovery",
		Severity: SeverityWarning,
		Description: "`pg_ctl stop -m immediate` sends `SIGQUIT` to the postmaster. Server " +
			"processes drop connections, no checkpoint is taken, and buffered changes are " +
			"left in memory. Recovery on the next start has to replay every record since the " +
			"last checkpoint; if WAL is corrupt, lost, or on different storage, committed " +
			"transactions can be lost. Use `-m smart` (default) or `-m fast` so the server " +
			"issues a shutdown checkpoint and closes cleanly; reserve `immediate` for the " +
			"\"the node is on fire\" case and pair it with a tested PITR procedure.",
		Check: checkZC1800,
	})
}

func checkZC1800(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "pg_ctl" {
		return nil
	}
	hasStop := false
	hasImmediate := false
	for i, arg := range cmd.Arguments {
		v := arg.String()
		if v == "stop" || v == "restart" {
			hasStop = true
		}
		if v == "-m" && i+1 < len(cmd.Arguments) {
			if cmd.Arguments[i+1].String() == "immediate" {
				hasImmediate = true
			}
		}
		if strings.HasPrefix(v, "-m") && len(v) > 2 && v[2:] == "immediate" {
			hasImmediate = true
		}
		if v == "--mode=immediate" {
			hasImmediate = true
		}
	}
	if !hasStop || !hasImmediate {
		return nil
	}
	return []Violation{{
		KataID: "ZC1800",
		Message: "`pg_ctl stop -m immediate` kills the postmaster without a shutdown " +
			"checkpoint — WAL replay on restart can lose committed transactions " +
			"if WAL is corrupt. Use `-m smart` or `-m fast` for routine shutdowns.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
