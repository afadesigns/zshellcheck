// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1945",
		Title:    "Warn on `bpftrace -e` / `bpftool prog load` — loads in-kernel eBPF from a script",
		Severity: SeverityWarning,
		Description: "`bpftrace -e '…'` compiles an inline script into an eBPF program and attaches " +
			"to kprobes, tracepoints, or uprobes; `bpftool prog load FILE pinned /sys/fs/bpf/…` " +
			"installs a pre-built program. Both require `CAP_BPF`/`CAP_SYS_ADMIN` and can read " +
			"arbitrary kernel/userland memory — every command a sibling process runs, every " +
			"syscall argument, every TCP payload. Pin the loaded program to a directory the " +
			"operator owns, gate invocation behind a runbook, and prefer a short-lived " +
			"`bpftrace -c CMD` window over long-running traces left on the host.",
		Check: checkZC1945,
	})
}

func checkZC1945(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "bpftrace":
		for _, arg := range cmd.Arguments {
			if arg.String() == "-e" {
				return zc1945Hit(cmd, "bpftrace -e")
			}
		}
	case "bpftool":
		if len(cmd.Arguments) >= 2 &&
			cmd.Arguments[0].String() == "prog" &&
			cmd.Arguments[1].String() == "load" {
			return zc1945Hit(cmd, "bpftool prog load")
		}
	}
	return nil
}

func zc1945Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1945",
		Message: "`" + form + "` loads an in-kernel eBPF program that can read arbitrary " +
			"kernel/userland memory — every syscall arg, every TCP payload. Gate behind a " +
			"runbook and prefer a short-lived `bpftrace -c CMD` over a pinned trace.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
