// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1643(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — $(<file)",
			input:    `echo "$(<file)"`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — bare cat file (logging, not capture)",
			input:    `cat /etc/os-release`,
			expected: []katas.Violation{},
		},
		{
			name:  `invalid — echo "$(cat /etc/hostname)"`,
			input: `echo "$(cat /etc/hostname)"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1643",
					Message: "`$(cat FILE)` forks cat just to read a file — use `$(<FILE)` (shell builtin, no fork).",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  `invalid — print -r -- "$(cat file)"`,
			input: `print -r -- "$(cat file)"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1643",
					Message: "`$(cat FILE)` forks cat just to read a file — use `$(<FILE)` (shell builtin, no fork).",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1643")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
