package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1360(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — ls -l (not sort-by-size)",
			input:    `ls -l`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — ls -S",
			input: `ls -S`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1360",
					Message: "Use Zsh `*(OL)` (largest-first) or `*(oL)` (smallest-first) glob qualifier instead of `ls -S`. No external process needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — ls -lS",
			input: `ls -lS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1360",
					Message: "Use Zsh `*(OL)` (largest-first) or `*(oL)` (smallest-first) glob qualifier instead of `ls -S`. No external process needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1360")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
