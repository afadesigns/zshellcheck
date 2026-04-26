// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1665",
		Title:    "Warn on `chrt -r` / `-f` — real-time scheduling class from a shell script",
		Severity: SeverityWarning,
		Description: "`chrt -r PRIO CMD` (SCHED_RR) and `chrt -f PRIO CMD` (SCHED_FIFO) launch " +
			"the child under a POSIX real-time scheduling class. An RT thread preempts " +
			"every normal-priority task until it voluntarily yields; a busy-loop or a " +
			"deadlock leaves the kernel with kworker, ksoftirqd, and sshd starved, often " +
			"forcing a hard reboot. Unless the binary is known-bounded (audio glitch-free " +
			"path, protocol timing loop), keep scripts on SCHED_OTHER — use `nice -n -5` or " +
			"a systemd unit with `CPUWeight=` / `IOWeight=` instead of `chrt -r`.",
		Check: checkZC1665,
	})
}

func checkZC1665(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "chrt" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-r" || v == "-f" || v == "--rr" || v == "--fifo" {
			return []Violation{{
				KataID: "ZC1665",
				Message: "`chrt " + v + "` puts the child on a real-time scheduling class — a " +
					"busy-loop or deadlock then starves kworker / sshd. Prefer `nice -n -5` " +
					"or a systemd unit with `CPUWeight=`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
