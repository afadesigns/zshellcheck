package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1287(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid cat with file",
			input:    `cat file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid cat with -n flag",
			input:    `cat -n file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid cat -v for visible chars",
			input: `cat -v file.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1287",
					Message: "Use Zsh `${(V)var}` to make control characters visible instead of `cat -v`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid cat -A for all visible",
			input: `cat -A file.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1287",
					Message: "Use Zsh `${(V)var}` to make control characters visible instead of `cat -v`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1287")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
