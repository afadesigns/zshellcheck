package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1091(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid string comparison",
			input:    `[[ $a == $b ]]`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid file test",
			input:    `[[ -f $file ]]`,
			expected: []katas.Violation{},
		},
		{
			name:     "invalid arithmetic -eq",
			input:    `[[ $a -eq 1 ]]`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1091",
					Message: "Use `(( ... ))` for arithmetic comparisons. e.g. `(( a < b ))` instead of `[[ a -lt b ]]`.",
					Line:    1,
					Column:  7, // -eq token column
				},
			},
		},
		{
			name:     "invalid arithmetic -lt",
			input:    `[[ $a -lt 5 ]]`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1091",
					Message: "Use `(( ... ))` for arithmetic comparisons. e.g. `(( a < b ))` instead of `[[ a -lt b ]]`.",
					Line:    1,
					Column:  7,
				},
			},
		},
		{
			name:     "invalid nested arithmetic",
			input:    `[[ $a -gt 0 && $b -lt 10 ]]`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1091",
					Message: "Use `(( ... ))` for arithmetic comparisons. e.g. `(( a < b ))` instead of `[[ a -lt b ]]`.",
					Line:    1,
					Column:  7, // -gt
				},
				{
					KataID:  "ZC1091",
					Message: "Use `(( ... ))` for arithmetic comparisons. e.g. `(( a < b ))` instead of `[[ a -lt b ]]`.",
					Line:    1,
					Column:  19, // -lt
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1091")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
