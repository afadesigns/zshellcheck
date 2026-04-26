// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1901(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt POSIX_BUILTINS` (explicit default)",
			input:    `unsetopt POSIX_BUILTINS`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NOMATCH` (unrelated)",
			input:    `setopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt POSIX_BUILTINS`",
			input: `setopt POSIX_BUILTINS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1901",
					Message: "`setopt POSIX_BUILTINS` switches Zsh to POSIX special-builtin rules — assignments before `export`/`readonly`/`eval` stop being local, silently leaking state. Scope any POSIX block with `emulate -LR sh` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_POSIX_BUILTINS`",
			input: `unsetopt NO_POSIX_BUILTINS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1901",
					Message: "`unsetopt NO_POSIX_BUILTINS` switches Zsh to POSIX special-builtin rules — assignments before `export`/`readonly`/`eval` stop being local, silently leaking state. Scope any POSIX block with `emulate -LR sh` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1901")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
