// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1775(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `timeout -k 5 30 cmd` (escalation configured)",
			input:    `timeout -k 5 30 cmd`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `timeout --kill-after=5 30 cmd`",
			input:    `timeout --kill-after=5 30 cmd`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `timeout` alone (no command, probably listing help)",
			input:    `timeout`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `timeout 30 cmd` (no escalation)",
			input: `timeout 30 cmd`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1775",
					Message: "`timeout` without `--kill-after` / `-k` only sends `SIGTERM` — a child that blocks or ignores it hangs the pipeline past the deadline. Add `--kill-after=N` so timeout escalates to `SIGKILL`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `timeout 5m long-running-job`",
			input: `timeout 5m long-running-job`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1775",
					Message: "`timeout` without `--kill-after` / `-k` only sends `SIGTERM` — a child that blocks or ignores it hangs the pipeline past the deadline. Add `--kill-after=N` so timeout escalates to `SIGKILL`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1775")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
