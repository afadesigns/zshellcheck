// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1973(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt POSIX_IDENTIFIERS` (restores default)",
			input:    `unsetopt POSIX_IDENTIFIERS`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NO_POSIX_IDENTIFIERS` (restores default)",
			input:    `setopt NO_POSIX_IDENTIFIERS`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt POSIX_IDENTIFIERS`",
			input: `setopt POSIX_IDENTIFIERS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1973",
					Message: "`setopt POSIX_IDENTIFIERS` restricts parameter names to ASCII; later `${café}`/`${π}` fail to parse and i18n-named libs stop loading. Scope with `emulate -LR sh` inside the helper instead of flipping globally.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_POSIX_IDENTIFIERS`",
			input: `unsetopt NO_POSIX_IDENTIFIERS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1973",
					Message: "`unsetopt NO_POSIX_IDENTIFIERS` restricts parameter names to ASCII; later `${café}`/`${π}` fail to parse and i18n-named libs stop loading. Scope with `emulate -LR sh` inside the helper instead of flipping globally.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1973")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
