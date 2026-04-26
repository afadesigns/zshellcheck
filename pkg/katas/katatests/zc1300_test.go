// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1300(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid ZSH_VERSION usage",
			input:    `echo $ZSH_VERSION`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid BASH_VERSINFO usage",
			input: `echo $BASH_VERSINFO`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1300",
					Message: "Avoid Bash version variables in Zsh — use `$ZSH_VERSION` instead. Bash version variables are undefined in Zsh.",
					Line:    1,
					Column:  6,
				},
			},
		},
		{
			name:  "invalid BASH_VERSION usage",
			input: `echo $BASH_VERSION`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1300",
					Message: "Avoid Bash version variables in Zsh — use `$ZSH_VERSION` instead. Bash version variables are undefined in Zsh.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1300")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
