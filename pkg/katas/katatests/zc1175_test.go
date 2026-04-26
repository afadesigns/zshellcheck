// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1175(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid tput cols",
			input:    `tput cols`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid tput setaf",
			input: `tput setaf 1`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1175",
					Message: "Use Zsh `%F{color}` / `%f` or `$fg[color]` / `$reset_color` instead of `tput`. Zsh handles ANSI colors natively without spawning external processes.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid tput sgr0",
			input: `tput sgr0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1175",
					Message: "Use Zsh `%F{color}` / `%f` or `$fg[color]` / `$reset_color` instead of `tput`. Zsh handles ANSI colors natively without spawning external processes.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1175")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
