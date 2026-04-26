// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1660(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — printf %d without width",
			input:    `printf '%d' 5`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — printf %-20s left-aligned string (space fill)",
			input:    `printf '%-20s' "$name"`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — printf %05d",
			input: `printf '%05d' $n`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1660",
					Message: "`printf '%0Nd'` forks for zero-padding — prefer Zsh `${(l:N::0:)n}` parameter-expansion pad (same for `(r:N::0:)` on the right).",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — printf %03d literal",
			input: `printf '%03d' 7`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1660",
					Message: "`printf '%0Nd'` forks for zero-padding — prefer Zsh `${(l:N::0:)n}` parameter-expansion pad (same for `(r:N::0:)` on the right).",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1660")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
