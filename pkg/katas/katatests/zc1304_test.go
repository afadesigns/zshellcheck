// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1304(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid ZSH_SUBSHELL usage",
			input:    `echo $ZSH_SUBSHELL`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid BASH_SUBSHELL usage",
			input: `echo $BASH_SUBSHELL`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1304",
					Message: "Avoid `$BASH_SUBSHELL` in Zsh — use `$ZSH_SUBSHELL` instead. `BASH_SUBSHELL` is Bash-specific.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1304")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
