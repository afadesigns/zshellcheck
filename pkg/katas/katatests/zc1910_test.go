// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1910(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt GLOB_STAR_SHORT` (explicit default)",
			input:    `unsetopt GLOB_STAR_SHORT`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt EXTENDED_GLOB` (unrelated)",
			input:    `setopt EXTENDED_GLOB`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt GLOB_STAR_SHORT`",
			input: `setopt GLOB_STAR_SHORT`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1910",
					Message: "`setopt GLOB_STAR_SHORT` turns bare `**` into `**/*` — `rm **` now wipes the tree. Keep the option off and spell `**/*` when recursion is actually wanted.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_GLOB_STAR_SHORT`",
			input: `unsetopt NO_GLOB_STAR_SHORT`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1910",
					Message: "`unsetopt NO_GLOB_STAR_SHORT` turns bare `**` into `**/*` — `rm **` now wipes the tree. Keep the option off and spell `**/*` when recursion is actually wanted.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1910")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
