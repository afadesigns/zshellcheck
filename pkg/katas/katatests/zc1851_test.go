package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1851(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `setopt FUNCTION_ARGZERO` (explicit default)",
			input:    `setopt FUNCTION_ARGZERO`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `unsetopt NOMATCH` (unrelated)",
			input:    `unsetopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `unsetopt FUNCTION_ARGZERO`",
			input: `unsetopt FUNCTION_ARGZERO`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1851",
					Message: "`unsetopt FUNCTION_ARGZERO` makes `$0` inside functions point at the outer script — breaks `log \"$0: ...\"` helpers and `case $0` dispatchers. Keep the option on; reach the script name explicitly via `$ZSH_ARGZERO` when needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `setopt NO_FUNCTION_ARGZERO`",
			input: `setopt NO_FUNCTION_ARGZERO`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1851",
					Message: "`setopt NO_FUNCTION_ARGZERO` makes `$0` inside functions point at the outer script — breaks `log \"$0: ...\"` helpers and `case $0` dispatchers. Keep the option on; reach the script name explicitly via `$ZSH_ARGZERO` when needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1851")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
