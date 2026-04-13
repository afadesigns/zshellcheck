package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1292(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid tr with character ranges",
			input:    `tr 'a-z' 'A-Z'`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid tr with POSIX classes",
			input:    `tr '[:upper:]' '[:lower:]'`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid tr with single char translation",
			input: `tr '/' '_'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1292",
					Message: "Use Zsh `${var////_}` for character substitution instead of `tr`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1292")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
