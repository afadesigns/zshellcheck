package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1284(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid Zsh split expansion",
			input:    `echo ${(s/:/)PATH}`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid cut with dot delimiter (covered by ZC1280)",
			input:    `cut -d. -f2`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid cut with colon delimiter",
			input: `cut -d: -f1`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1284",
					Message: "Use Zsh parameter expansion `${(s:sep:)var}` for field splitting instead of `cut -d -f`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid cut with comma delimiter",
			input: `cut -d, -f3`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1284",
					Message: "Use Zsh parameter expansion `${(s:sep:)var}` for field splitting instead of `cut -d -f`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1284")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
