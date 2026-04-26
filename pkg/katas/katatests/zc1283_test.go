// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1283(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid setopt usage",
			input:    `setopt noglob`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid unsetopt usage",
			input:    `unsetopt noglob`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid set without -o",
			input:    `set -e`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid set -o",
			input: `set -o noglob`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1283",
					Message: "Use `setopt` instead of `set -o` in Zsh scripts. `setopt` is the native Zsh idiom.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1283")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
