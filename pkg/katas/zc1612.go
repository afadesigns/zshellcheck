// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

var zc1612HardeningDisables = map[string]string{
	"kernel.yama.ptrace_scope=0":         "YAMA ptrace scope (lets any process attach)",
	"kernel.kptr_restrict=0":             "kernel pointer restriction (leaks kptrs to /proc)",
	"kernel.dmesg_restrict=0":            "dmesg restriction (unprivileged users read ring buffer)",
	"kernel.unprivileged_bpf_disabled=0": "unprivileged BPF gate (any user loads BPF)",
	"net.core.bpf_jit_harden=0":          "BPF JIT hardening (JIT-spray mitigations off)",
	"kernel.perf_event_paranoid=-1":      "perf_event paranoid (unprivileged perf access)",
	"kernel.perf_event_paranoid=0":       "perf_event paranoid (unprivileged perf access)",
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1612",
		Title:    "Warn on `sysctl -w` disabling kernel hardening knobs",
		Severity: SeverityWarning,
		Description: "Several sysctl knobs exist specifically to constrain what unprivileged " +
			"users can do — `kernel.yama.ptrace_scope`, `kernel.kptr_restrict`, " +
			"`kernel.dmesg_restrict`, `kernel.unprivileged_bpf_disabled`, " +
			"`net.core.bpf_jit_harden`, and `kernel.perf_event_paranoid`. Setting any of them " +
			"to the lowest-restriction value removes a distinct defense-in-depth layer: " +
			"unrelated processes can ptrace each other, kernel pointers leak to `/proc`, " +
			"unprivileged users read kernel ring buffers, BPF JIT-spray mitigations disappear. " +
			"Leave these defaults alone unless a measured performance or debugging need justifies it.",
		Check: checkZC1612,
	})
}

func checkZC1612(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "sysctl" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if note, ok := zc1612HardeningDisables[v]; ok {
			return []Violation{{
				KataID: "ZC1612",
				Message: "`sysctl ... " + v + "` disables " + note + " — defense-in-depth " +
					"loss. Leave the default unless a measured need justifies it.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
