// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1869(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt RC_EXPAND_PARAM` (explicit default)",
			input:    `unsetopt RC_EXPAND_PARAM`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NOMATCH` (unrelated)",
			input:    `setopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt RC_EXPAND_PARAM`",
			input: `setopt RC_EXPAND_PARAM`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1869",
					Message: "`setopt RC_EXPAND_PARAM` distributes literal prefix/suffix across every array element — `cp src/${arr[@]}.bak dst` silently rewrites as `cp src/a.bak src/b.bak dst`. Keep it off; opt in per-use with `${^arr}`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_RC_EXPAND_PARAM`",
			input: `unsetopt NO_RC_EXPAND_PARAM`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1869",
					Message: "`unsetopt NO_RC_EXPAND_PARAM` distributes literal prefix/suffix across every array element — `cp src/${arr[@]}.bak dst` silently rewrites as `cp src/a.bak src/b.bak dst`. Keep it off; opt in per-use with `${^arr}`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1869")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
