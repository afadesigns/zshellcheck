package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1967(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt PROMPT_SUBST`",
			input:    `unsetopt PROMPT_SUBST`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NO_PROMPT_SUBST`",
			input:    `setopt NO_PROMPT_SUBST`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt PROMPT_SUBST`",
			input: `setopt PROMPT_SUBST`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1967",
					Message: "`setopt PROMPT_SUBST` re-runs command substitution on every prompt redraw — a branch/host/dir value with `$(…)` executes each render. Prefer `%n`/`%d`/`%~`/`vcs_info`, or scope via `LOCAL_OPTIONS`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_PROMPT_SUBST`",
			input: `unsetopt NO_PROMPT_SUBST`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1967",
					Message: "`unsetopt NO_PROMPT_SUBST` re-runs command substitution on every prompt redraw — a branch/host/dir value with `$(…)` executes each render. Prefer `%n`/`%d`/`%~`/`vcs_info`, or scope via `LOCAL_OPTIONS`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1967")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
