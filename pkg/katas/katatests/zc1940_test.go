// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1940(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt POSIX_ARGZERO` (explicit default)",
			input:    `unsetopt POSIX_ARGZERO`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt FUNCTION_ARGZERO` (different option)",
			input:    `setopt FUNCTION_ARGZERO`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt POSIX_ARGZERO`",
			input: `setopt POSIX_ARGZERO`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1940",
					Message: "`setopt POSIX_ARGZERO` freezes `$0` to the outer script name — loggers and `case $0` dispatch inside functions lose call-site context. Scope with `emulate -LR sh` instead of flipping globally.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_POSIX_ARGZERO`",
			input: `unsetopt NO_POSIX_ARGZERO`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1940",
					Message: "`unsetopt NO_POSIX_ARGZERO` freezes `$0` to the outer script name — loggers and `case $0` dispatch inside functions lose call-site context. Scope with `emulate -LR sh` instead of flipping globally.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1940")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
