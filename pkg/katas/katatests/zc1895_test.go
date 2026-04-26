// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1895(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt NUMERIC_GLOB_SORT` (explicit default)",
			input:    `unsetopt NUMERIC_GLOB_SORT`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NOMATCH` (unrelated)",
			input:    `setopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt NUMERIC_GLOB_SORT`",
			input: `setopt NUMERIC_GLOB_SORT`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1895",
					Message: "`setopt NUMERIC_GLOB_SORT` switches every later glob to numeric sort — log rotations sorted on numeric suffixes silently shuffle. Keep it off; use the per-glob `*(n)` qualifier when needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_NUMERIC_GLOB_SORT`",
			input: `unsetopt NO_NUMERIC_GLOB_SORT`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1895",
					Message: "`unsetopt NO_NUMERIC_GLOB_SORT` switches every later glob to numeric sort — log rotations sorted on numeric suffixes silently shuffle. Keep it off; use the per-glob `*(n)` qualifier when needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1895")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
