// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1168(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:  "invalid readarray",
			input: `readarray -t arr`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1168",
					Message: "Use Zsh `${(f)$(cmd)}` instead of `readarray`. `readarray`/`mapfile` are Bash builtins not available in Zsh.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid mapfile",
			input: `mapfile -t lines`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1168",
					Message: "Use Zsh `${(f)$(cmd)}` instead of `mapfile`. `readarray`/`mapfile` are Bash builtins not available in Zsh.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1168")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
