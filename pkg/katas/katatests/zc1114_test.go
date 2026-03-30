package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1114(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid mktemp -d",
			input:    `mktemp -d`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid mktemp",
			input: `mktemp /tmp/zsh.XXXXXX`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1114",
					Message: "Consider using Zsh `=(cmd)` for temporary files instead of `mktemp`. Zsh auto-cleans temporary files created with `=(...)` process substitution.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid bare mktemp",
			input: `mktemp -t prefix`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1114",
					Message: "Consider using Zsh `=(cmd)` for temporary files instead of `mktemp`. Zsh auto-cleans temporary files created with `=(...)` process substitution.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1114")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
