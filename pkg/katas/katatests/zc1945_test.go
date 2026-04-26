// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1945(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `bpftrace -l 'tracepoint:syscalls:*'` (list only)",
			input:    `bpftrace -l tracepoint:syscalls:*`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `bpftool prog show`",
			input:    `bpftool prog show`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `bpftrace -e 'tracepoint:syscalls:sys_enter_openat{printf(...)}'`",
			input: `bpftrace -e 'tracepoint:syscalls:sys_enter_openat{printf(...)}'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1945",
					Message: "`bpftrace -e` loads an in-kernel eBPF program that can read arbitrary kernel/userland memory — every syscall arg, every TCP payload. Gate behind a runbook and prefer a short-lived `bpftrace -c CMD` over a pinned trace.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `bpftool prog load prog.o /sys/fs/bpf/spy`",
			input: `bpftool prog load prog.o /sys/fs/bpf/spy`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1945",
					Message: "`bpftool prog load` loads an in-kernel eBPF program that can read arbitrary kernel/userland memory — every syscall arg, every TCP payload. Gate behind a runbook and prefer a short-lived `bpftrace -c CMD` over a pinned trace.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1945")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
