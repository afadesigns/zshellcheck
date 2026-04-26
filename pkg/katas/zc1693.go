// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1693",
		Title:    "Warn on `ionice -c 1` — real-time I/O class starves every other disk consumer",
		Severity: SeverityWarning,
		Description: "`ionice -c 1` (real-time I/O scheduling class) promotes the child above " +
			"every best-effort (class 2) and idle (class 3) task queued against the same " +
			"device. A busy workload — `rsync`, `dd`, database backup — then blocks sshd " +
			"reads, systemd journal writes, and every other process until it yields, which " +
			"for sequential I/O is effectively never. If the intent is \"fast I/O\", stay on " +
			"class 2 and let CFQ / BFQ handle it; reserve class 1 for latency-critical " +
			"paths launched by a scheduler that knows how to cap duration.",
		Check: checkZC1693,
	})
}

func checkZC1693(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "ionice" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-c1" {
			return zc1693Hit(cmd)
		}
		if v == "-c" && i+1 < len(cmd.Arguments) && cmd.Arguments[i+1].String() == "1" {
			return zc1693Hit(cmd)
		}
	}
	return nil
}

func zc1693Hit(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1693",
		Message: "`ionice -c 1` puts the child in the real-time I/O class — a long-running " +
			"workload starves sshd / journald / the rest of the host. Stay on class 2.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
