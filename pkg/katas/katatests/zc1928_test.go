// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1928(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt SHARE_HISTORY` (explicit default)",
			input:    `unsetopt SHARE_HISTORY`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt INC_APPEND_HISTORY` (safer alternative)",
			input:    `setopt INC_APPEND_HISTORY`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt SHARE_HISTORY`",
			input: `setopt SHARE_HISTORY`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1928",
					Message: "`setopt SHARE_HISTORY` flushes every command into every sibling zsh session — secrets typed in one terminal surface in `fc -l` of every other. Prefer `setopt INC_APPEND_HISTORY` plus `HIST_IGNORE_SPACE` for safer isolation.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_SHARE_HISTORY`",
			input: `unsetopt NO_SHARE_HISTORY`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1928",
					Message: "`unsetopt NO_SHARE_HISTORY` flushes every command into every sibling zsh session — secrets typed in one terminal surface in `fc -l` of every other. Prefer `setopt INC_APPEND_HISTORY` plus `HIST_IGNORE_SPACE` for safer isolation.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1928")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
