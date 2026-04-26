// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1481(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — unset TMPDIR",
			input:    `unset TMPDIR`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — export HISTFILE=~/.zsh_history",
			input:    `export HISTFILE=~/.zsh_history`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — unset HISTFILE",
			input: `unset HISTFILE`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1481",
					Message: "`unset HISTFILE` disables shell history — textbook post-compromise tactic. Legitimate alternative: `HISTCONTROL=ignorespace` plus leading-space prefix.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — export HISTFILE=/dev/null",
			input: `export HISTFILE=/dev/null`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1481",
					Message: "`HISTFILE=/dev/null` disables shell history — textbook post-compromise tactic. Legitimate alternative: `HISTCONTROL=ignorespace` plus leading-space prefix.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — export HISTSIZE=0",
			input: `export HISTSIZE=0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1481",
					Message: "`HISTSIZE=0` disables shell history — textbook post-compromise tactic. Legitimate alternative: `HISTCONTROL=ignorespace` plus leading-space prefix.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1481")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
