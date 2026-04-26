// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC2003(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt KSH_ZERO_SUBSCRIPT` (default)",
			input:    `unsetopt KSH_ZERO_SUBSCRIPT`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NO_KSH_ZERO_SUBSCRIPT`",
			input:    `setopt NO_KSH_ZERO_SUBSCRIPT`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt KSH_ZERO_SUBSCRIPT`",
			input: `setopt KSH_ZERO_SUBSCRIPT`,
			expected: []katas.Violation{
				{
					KataID:  "ZC2003",
					Message: "`setopt KSH_ZERO_SUBSCRIPT` stops aliasing `$arr[0]` to `$arr[1]` — every later read of `$arr[0]` silently returns empty and `arr[0]=new` stops updating the first element. Use `$arr[1]` explicitly.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_KSH_ZERO_SUBSCRIPT`",
			input: `unsetopt NO_KSH_ZERO_SUBSCRIPT`,
			expected: []katas.Violation{
				{
					KataID:  "ZC2003",
					Message: "`unsetopt NO_KSH_ZERO_SUBSCRIPT` stops aliasing `$arr[0]` to `$arr[1]` — every later read of `$arr[0]` silently returns empty and `arr[0]=new` stops updating the first element. Use `$arr[1]` explicitly.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC2003")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
