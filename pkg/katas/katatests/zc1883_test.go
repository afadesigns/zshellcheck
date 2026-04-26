// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1883(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt PATH_SCRIPT` (explicit default)",
			input:    `unsetopt PATH_SCRIPT`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NOMATCH` (unrelated)",
			input:    `setopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt PATH_SCRIPT`",
			input: `setopt PATH_SCRIPT`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1883",
					Message: "`setopt PATH_SCRIPT` lets `.`/`source` fall back to `$PATH` when a literal path misses — a dropper in `~/bin` or `./` runs inside the current shell with every exported secret. Keep the option off; always use explicit paths.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_PATH_SCRIPT`",
			input: `unsetopt NO_PATH_SCRIPT`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1883",
					Message: "`unsetopt NO_PATH_SCRIPT` lets `.`/`source` fall back to `$PATH` when a literal path misses — a dropper in `~/bin` or `./` runs inside the current shell with every exported secret. Keep the option off; always use explicit paths.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1883")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
