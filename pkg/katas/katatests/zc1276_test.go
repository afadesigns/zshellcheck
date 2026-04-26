// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1276(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid brace expansion",
			input:    `echo {1..10}`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid other command",
			input:    `echo hello`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid seq usage",
			input: `seq 1 10`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1276",
					Message: "Use Zsh brace expansion `{start..end}` instead of `seq`. Brace expansion is built-in and avoids forking.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid seq with single arg",
			input: `seq 5`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1276",
					Message: "Use Zsh brace expansion `{start..end}` instead of `seq`. Brace expansion is built-in and avoids forking.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1276")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
