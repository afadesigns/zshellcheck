package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1327(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid fc usage",
			input:    `fc -l`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid history without flags",
			input:    `history 10`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid history -c usage",
			input: `history -c`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1327",
					Message: "Avoid `history -c` in Zsh — Bash history flags differ. Use `fc` commands for Zsh history management.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1327")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
