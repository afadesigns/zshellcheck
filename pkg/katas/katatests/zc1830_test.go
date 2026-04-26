// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1830(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `setopt NOMATCH`",
			input:    `setopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `unsetopt HIST_IGNORE_DUPS` (unrelated)",
			input:    `unsetopt HIST_IGNORE_DUPS`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `unsetopt NOMATCH`",
			input: `unsetopt NOMATCH`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1830",
					Message: "`unsetopt NOMATCH` silences Zsh's unmatched-glob error — typos pass through literally. Use `*(N)` per-glob or scope inside a function with `setopt LOCAL_OPTIONS; setopt NULL_GLOB`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `setopt NO_NOMATCH`",
			input: `setopt NO_NOMATCH`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1830",
					Message: "`setopt NO_NOMATCH` silences Zsh's unmatched-glob error — typos pass through literally. Use `*(N)` per-glob or scope inside a function with `setopt LOCAL_OPTIONS; setopt NULL_GLOB`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1830")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
