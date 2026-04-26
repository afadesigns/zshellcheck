// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1790(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `setopt PIPE_FAIL` (enabling)",
			input:    `setopt PIPE_FAIL`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `unsetopt NOMATCH` (unrelated)",
			input:    `unsetopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `unsetopt PIPE_FAIL`",
			input: `unsetopt PIPE_FAIL`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1790",
					Message: "`unsetopt PIPE_FAIL` returns the shell to last-command-only pipeline exit — `cmd1 | cmd2` now ignores `cmd1` failures. Scope the change to a subshell or function with `emulate -L zsh` instead of flipping it globally.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `setopt NOPIPEFAIL`",
			input: `setopt NOPIPEFAIL`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1790",
					Message: "`setopt NOPIPEFAIL` returns the shell to last-command-only pipeline exit — `cmd1 | cmd2` now ignores `cmd1` failures. Scope the change to a subshell or function with `emulate -L zsh` instead of flipping it globally.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1790")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
