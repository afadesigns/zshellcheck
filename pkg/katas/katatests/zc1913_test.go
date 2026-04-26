// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1913(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt ALIAS_FUNC_DEF` (explicit default)",
			input:    `unsetopt ALIAS_FUNC_DEF`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt AUTO_CD` (unrelated)",
			input:    `setopt AUTO_CD`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt ALIAS_FUNC_DEF`",
			input: `setopt ALIAS_FUNC_DEF`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1913",
					Message: "`setopt ALIAS_FUNC_DEF` lets a function silently shadow an alias — one sourced rc file replaces your function with the alias, no error surfaces. Keep it off; quote the name if the override is intentional.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_ALIAS_FUNC_DEF`",
			input: `unsetopt NO_ALIAS_FUNC_DEF`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1913",
					Message: "`unsetopt NO_ALIAS_FUNC_DEF` lets a function silently shadow an alias — one sourced rc file replaces your function with the alias, no error surfaces. Keep it off; quote the name if the override is intentional.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1913")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
