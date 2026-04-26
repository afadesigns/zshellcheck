// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1887(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt POSIX_TRAPS` (explicit default)",
			input:    `unsetopt POSIX_TRAPS`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NOMATCH` (unrelated)",
			input:    `setopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt POSIX_TRAPS`",
			input: `setopt POSIX_TRAPS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1887",
					Message: "`setopt POSIX_TRAPS` flips `trap ... EXIT` inside functions from function-return to shell-exit scope — per-call cleanup leaks across the whole shell, TRAPZERR helpers stop firing. Keep the option off.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_POSIX_TRAPS`",
			input: `unsetopt NO_POSIX_TRAPS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1887",
					Message: "`unsetopt NO_POSIX_TRAPS` flips `trap ... EXIT` inside functions from function-return to shell-exit scope — per-call cleanup leaks across the whole shell, TRAPZERR helpers stop firing. Keep the option off.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1887")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
