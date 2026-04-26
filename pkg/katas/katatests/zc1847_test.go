// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1847(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt CHASE_LINKS` (explicit default)",
			input:    `unsetopt CHASE_LINKS`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NOMATCH` (unrelated)",
			input:    `setopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt CHASE_LINKS`",
			input: `setopt CHASE_LINKS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1847",
					Message: "`setopt CHASE_LINKS` makes every `cd` resolve symlinks to the physical inode — `cd releases/current` lands in the release dir, breaking `..` navigation. Keep it off; use `cd -P target` one-shot when needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_CHASE_LINKS`",
			input: `unsetopt NO_CHASE_LINKS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1847",
					Message: "`unsetopt NO_CHASE_LINKS` makes every `cd` resolve symlinks to the physical inode — `cd releases/current` lands in the release dir, breaking `..` navigation. Keep it off; use `cd -P target` one-shot when needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1847")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
