package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1271(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid command -v usage",
			input:    `command -v git`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid type command",
			input:    `type git`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid which usage",
			input: `which git`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1271",
					Message: "Use `command -v` instead of `which`. `command -v` is POSIX-compliant and built into Zsh.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1271")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
