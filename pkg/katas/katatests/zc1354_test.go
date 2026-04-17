package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1354(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — type without flags",
			input:    `type grep`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — type -t",
			input: `type -t grep`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1354",
					Message: "Use Zsh `whence -w` (category), `whence -a` (all), or `whence -p` (path) instead of Bash-specific `type -t`/`-a`/`-P`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — type -P",
			input: `type -P grep`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1354",
					Message: "Use Zsh `whence -w` (category), `whence -a` (all), or `whence -p` (path) instead of Bash-specific `type -t`/`-a`/`-P`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1354")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
