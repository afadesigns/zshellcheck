// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1828(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `gcore --help`",
			input:    `gcore --help`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `strace ls` (trace a child, not ptrace-attach)",
			input:    `strace ls`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `gcore 1234`",
			input: `gcore 1234`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1828",
					Message: "`gcore PID` attaches via ptrace — target memory, env, and syscall args are exposed. Production scripts should not run ptrace; use `coredumpctl` on a captured core instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `strace -f -p 1234`",
			input: `strace -f -p 1234`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1828",
					Message: "`strace -p PID` attaches via ptrace — target memory, env, and syscall args are exposed. Production scripts should not run ptrace; use `coredumpctl` on a captured core instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1828")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
