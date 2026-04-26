// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1108(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid tr delete",
			input:    `tr -d '\n'`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid tr squeeze",
			input:    `tr -s ' '`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid tr with different sets",
			input:    `tr ':' '\n'`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid tr lowercase to uppercase POSIX",
			input: `tr '[:lower:]' '[:upper:]'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1108",
					Message: "Use `${(U)var}` for case conversion instead of `tr`. Zsh parameter expansion flags avoid spawning an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid tr uppercase to lowercase range",
			input: `tr 'A-Z' 'a-z'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1108",
					Message: "Use `${(L)var}` for case conversion instead of `tr`. Zsh parameter expansion flags avoid spawning an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1108")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
