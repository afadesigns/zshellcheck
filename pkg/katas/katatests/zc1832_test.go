// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1832(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `limit coredumpsize 0` (disable cores)",
			input:    `limit coredumpsize 0`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `limit stacksize unlimited` (unrelated resource)",
			input:    `limit stacksize unlimited`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `limit coredumpsize unlimited`",
			input: `limit coredumpsize unlimited`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1832",
					Message: "`limit coredumpsize unlimited` enables unbounded core dumps (Zsh-specific `limit` spelling of `ulimit -c unlimited`). A setuid crash drops its memory to disk as a world-readable file — leave the ceiling at the distro default.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unlimit coredumpsize`",
			input: `unlimit coredumpsize`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1832",
					Message: "`unlimit coredumpsize` enables unbounded core dumps (Zsh-specific `limit` spelling of `ulimit -c unlimited`). A setuid crash drops its memory to disk as a world-readable file — leave the ceiling at the distro default.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1832")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
