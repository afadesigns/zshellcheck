package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1290(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid sort -n with -r flag",
			input:    `sort -n -r file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid sort without -n",
			input:    `sort file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid sort -n alone",
			input: `sort -n file.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1290",
					Message: "Use Zsh `${(n)array}` for numeric sorting instead of `sort -n`. The `(n)` flag sorts numerically in-shell.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1290")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
