// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1897(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt SH_GLOB` (explicit default)",
			input:    `unsetopt SH_GLOB`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NOMATCH` (unrelated)",
			input:    `setopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt SH_GLOB`",
			input: `setopt SH_GLOB`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1897",
					Message: "`setopt SH_GLOB` disables Zsh-extended glob patterns — `*(N)` qualifiers, `<1-10>` ranges, and `(alt1|alt2)` alternation raise parse errors. Keep the option off; scope with `LOCAL_OPTIONS` if strict POSIX is needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_SH_GLOB`",
			input: `unsetopt NO_SH_GLOB`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1897",
					Message: "`unsetopt NO_SH_GLOB` disables Zsh-extended glob patterns — `*(N)` qualifiers, `<1-10>` ranges, and `(alt1|alt2)` alternation raise parse errors. Keep the option off; scope with `LOCAL_OPTIONS` if strict POSIX is needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1897")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
