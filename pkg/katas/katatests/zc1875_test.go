// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1875(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt RC_QUOTES` (explicit default)",
			input:    `unsetopt RC_QUOTES`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NOMATCH` (unrelated)",
			input:    `setopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt RC_QUOTES`",
			input: `setopt RC_QUOTES`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1875",
					Message: "`setopt RC_QUOTES` reinterprets `''` inside single quotes as a literal apostrophe — `'it''s'` flips from `its` to `it's`, breaking tokens and SQL. Use double quotes or `\\'` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_RC_QUOTES`",
			input: `unsetopt NO_RC_QUOTES`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1875",
					Message: "`unsetopt NO_RC_QUOTES` reinterprets `''` inside single quotes as a literal apostrophe — `'it''s'` flips from `its` to `it's`, breaking tokens and SQL. Use double quotes or `\\'` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1875")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
