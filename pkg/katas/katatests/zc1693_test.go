// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1693(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — ionice -c 2 (best-effort)",
			input:    `ionice -c 2 -n 4 cmd`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — ionice -c 3 (idle)",
			input:    `ionice -c 3 cmd`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — ionice -c 1 (real-time, split)",
			input: `ionice -c 1 -n 0 cmd`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1693",
					Message: "`ionice -c 1` puts the child in the real-time I/O class — a long-running workload starves sshd / journald / the rest of the host. Stay on class 2.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — ionice -c1 (real-time, joined)",
			input: `ionice -c1 cmd`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1693",
					Message: "`ionice -c 1` puts the child in the real-time I/O class — a long-running workload starves sshd / journald / the rest of the host. Stay on class 2.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1693")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
