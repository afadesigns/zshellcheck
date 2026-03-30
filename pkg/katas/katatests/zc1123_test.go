package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1123(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid uname -r",
			input:    `uname -r`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid uname -m",
			input:    `uname -m`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid uname -s",
			input: `uname -s`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1123",
					Message: "Use `$OSTYPE` instead of `uname -s` for OS detection. Zsh maintains `$OSTYPE` as a built-in variable.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1123")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
