package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1932(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `setopt GLOBAL_EXPORT` (explicit default)",
			input:    `setopt GLOBAL_EXPORT`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `unsetopt NOMATCH` (unrelated)",
			input:    `unsetopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `unsetopt GLOBAL_EXPORT`",
			input: `unsetopt GLOBAL_EXPORT`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1932",
					Message: "`unsetopt GLOBAL_EXPORT` makes `typeset -x` exports function-local — helper functions that set `PATH`/`VIRTUAL_ENV`/`AWS_*` no longer propagate to callers. Keep it on; scope temporary exports in a subshell instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `setopt NO_GLOBAL_EXPORT`",
			input: `setopt NO_GLOBAL_EXPORT`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1932",
					Message: "`setopt NO_GLOBAL_EXPORT` makes `typeset -x` exports function-local — helper functions that set `PATH`/`VIRTUAL_ENV`/`AWS_*` no longer propagate to callers. Keep it on; scope temporary exports in a subshell instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1932")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
