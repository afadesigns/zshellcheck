package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1715(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — Zsh-idiomatic `read \"var?prompt\"`",
			input:    `read "name?Enter your name: "`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `read -r line` (no -p)",
			input:    `read -r line`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `read var` (bare)",
			input:    `read var`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `read -p \"Prompt: \" name`",
			input: `read -p "Prompt: " name`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1715",
					Message: "`read -p` triggers Zsh's coprocess reader, not Bash's prompt — the variable stays empty. Use `read \"var?Prompt: \"` (the `?` after the variable name introduces the prompt).",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `read -rp \"Prompt: \" name` (combined short flags)",
			input: `read -rp "Prompt: " name`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1715",
					Message: "`read -rp` triggers Zsh's coprocess reader, not Bash's prompt — the variable stays empty. Use `read \"var?Prompt: \"` (the `?` after the variable name introduces the prompt).",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1715")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
