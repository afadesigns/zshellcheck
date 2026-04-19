package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1789(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `setopt EXTENDED_GLOB`",
			input:    `setopt EXTENDED_GLOB`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `unsetopt CORRECT` (turning off is fine)",
			input:    `unsetopt CORRECT`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt CORRECT`",
			input: `setopt CORRECT`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1789",
					Message: "`setopt CORRECT` enables `CORRECT` — Zsh spellcheck silently rewrites tokens that look mistyped. In a script that corrupts file paths and steals stdin for the correction prompt. Keep in `~/.zshrc`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `setopt CORRECT_ALL`",
			input: `setopt CORRECT_ALL`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1789",
					Message: "`setopt CORRECT_ALL` enables `CORRECT_ALL` — Zsh spellcheck silently rewrites tokens that look mistyped. In a script that corrupts file paths and steals stdin for the correction prompt. Keep in `~/.zshrc`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `set -o correctall`",
			input: `set -o correctall`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1789",
					Message: "`set -o correctall` enables `CORRECT_ALL` — Zsh spellcheck silently rewrites tokens that look mistyped. In a script that corrupts file paths and steals stdin for the correction prompt. Keep in `~/.zshrc`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1789")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
