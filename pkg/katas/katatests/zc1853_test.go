// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1853(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt MARK_DIRS` (explicit default)",
			input:    `unsetopt MARK_DIRS`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NOMATCH` (unrelated)",
			input:    `setopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt MARK_DIRS`",
			input: `setopt MARK_DIRS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1853",
					Message: "`setopt MARK_DIRS` appends a trailing `/` to every glob-matched directory — `[[ -f \"$f\" ]]` and `rm -f *` start skipping, hash maps keyed on basenames double up. Keep the option off; use `*(/)` when you need dirs only.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_MARK_DIRS`",
			input: `unsetopt NO_MARK_DIRS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1853",
					Message: "`unsetopt NO_MARK_DIRS` appends a trailing `/` to every glob-matched directory — `[[ -f \"$f\" ]]` and `rm -f *` start skipping, hash maps keyed on basenames double up. Keep the option off; use `*(/)` when you need dirs only.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1853")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
