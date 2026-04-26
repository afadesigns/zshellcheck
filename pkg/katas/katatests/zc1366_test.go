// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1366(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — limit builtin",
			input:    `limit cputime 10`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — ulimit",
			input: `ulimit -t 10`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1366",
					Message: "Use Zsh `limit` (human-readable) or `limit -s` (stdout-only) instead of POSIX `ulimit` for Zsh-native resource queries.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1366")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
