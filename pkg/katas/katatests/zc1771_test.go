// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1771(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `alias ll='ls -l'` (regular alias)",
			input:    `alias ll='ls -l'`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `alias` (no args, lists aliases)",
			input:    `alias`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `alias -g G='| grep'`",
			input: `alias -g G='| grep'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1771",
					Message: "`alias -g` defines a global alias that expands outside command position — a surprise for anyone reading the script later. Prefer a function, or keep global aliases in `~/.zshrc` where they are discoverable.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `alias -s log=less`",
			input: `alias -s log=less`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1771",
					Message: "`alias -s` defines a suffix alias that expands outside command position — a surprise for anyone reading the script later. Prefer a function, or keep suffix aliases in `~/.zshrc` where they are discoverable.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1771")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
