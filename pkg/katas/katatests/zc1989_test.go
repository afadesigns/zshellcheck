// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1989(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt REMATCH_PCRE` (default)",
			input:    `unsetopt REMATCH_PCRE`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NO_REMATCH_PCRE`",
			input:    `setopt NO_REMATCH_PCRE`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt REMATCH_PCRE`",
			input: `setopt REMATCH_PCRE`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1989",
					Message: "`setopt REMATCH_PCRE` swaps `[[ =~ ]]` from POSIX ERE to PCRE — `\\b`, `\\d`, lookahead, `(?i)` change meaning across every later match. Prefer `pcre_match` where PCRE is needed or scope with `LOCAL_OPTIONS`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_REMATCH_PCRE`",
			input: `unsetopt NO_REMATCH_PCRE`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1989",
					Message: "`unsetopt NO_REMATCH_PCRE` swaps `[[ =~ ]]` from POSIX ERE to PCRE — `\\b`, `\\d`, lookahead, `(?i)` change meaning across every later match. Prefer `pcre_match` where PCRE is needed or scope with `LOCAL_OPTIONS`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1989")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
