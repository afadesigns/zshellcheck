// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1675(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — plain export VAR=value",
			input:    `export VAR=value`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — export multiple plain names",
			input:    `export PATH HOME`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — export -f function",
			input: `export -f my_func`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1675",
					Message: "`export -f` is Bash-only — use `typeset -fx` in Zsh.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — export -n VAR",
			input: `export -n VAR`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1675",
					Message: "`export -n` is Bash-only — use `typeset +x` in Zsh.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1675")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
