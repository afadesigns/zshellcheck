// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1665(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — chrt -o 0 (SCHED_OTHER)",
			input:    `chrt -o 0 /usr/bin/cmd`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — chrt listing priority",
			input:    `chrt -p 1234`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — chrt -r 99 cmd",
			input: `chrt -r 99 /usr/bin/cmd`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1665",
					Message: "`chrt -r` puts the child on a real-time scheduling class — a busy-loop or deadlock then starves kworker / sshd. Prefer `nice -n -5` or a systemd unit with `CPUWeight=`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — chrt -f 50 cmd",
			input: `chrt -f 50 /usr/bin/cmd`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1665",
					Message: "`chrt -f` puts the child on a real-time scheduling class — a busy-loop or deadlock then starves kworker / sshd. Prefer `nice -n -5` or a systemd unit with `CPUWeight=`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1665")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
