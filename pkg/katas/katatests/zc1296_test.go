// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1296(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid setopt usage",
			input:    `setopt extendedglob`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid shopt usage",
			input: `shopt -s extglob`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1296",
					Message: "Avoid `shopt` in Zsh — it is a Bash builtin. Use `setopt`/`unsetopt` for Zsh shell options.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1296")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
