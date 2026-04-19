package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1838(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt GLOB_DOTS` (explicit default)",
			input:    `unsetopt GLOB_DOTS`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NOMATCH` (unrelated)",
			input:    `setopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt GLOB_DOTS`",
			input: `setopt GLOB_DOTS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1838",
					Message: "`setopt GLOB_DOTS` makes every bare `*` also match hidden files — `rm *` quietly destroys `.git/`, `cp -r *` copies `.env`. Keep the option alone; request dotfiles per-glob with `*(D)`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_GLOB_DOTS`",
			input: `unsetopt NO_GLOB_DOTS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1838",
					Message: "`unsetopt NO_GLOB_DOTS` makes every bare `*` also match hidden files — `rm *` quietly destroys `.git/`, `cp -r *` copies `.env`. Keep the option alone; request dotfiles per-glob with `*(D)`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1838")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
