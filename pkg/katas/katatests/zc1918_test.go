package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1918(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt HIST_SUBST_PATTERN` (explicit default)",
			input:    `unsetopt HIST_SUBST_PATTERN`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt EXTENDED_HISTORY` (unrelated)",
			input:    `setopt EXTENDED_HISTORY`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt HIST_SUBST_PATTERN`",
			input: `setopt HIST_SUBST_PATTERN`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1918",
					Message: "`setopt HIST_SUBST_PATTERN` switches `:s` history/param modifiers to pattern matching — literal `*`/`?`/`^` suddenly act as glob metacharacters. Keep it off; use `${var//pat/rep}` when you actually want pattern substitution.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_HIST_SUBST_PATTERN`",
			input: `unsetopt NO_HIST_SUBST_PATTERN`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1918",
					Message: "`unsetopt NO_HIST_SUBST_PATTERN` switches `:s` history/param modifiers to pattern matching — literal `*`/`?`/`^` suddenly act as glob metacharacters. Keep it off; use `${var//pat/rep}` when you actually want pattern substitution.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1918")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
