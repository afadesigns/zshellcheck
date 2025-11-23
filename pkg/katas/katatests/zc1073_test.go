package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1073(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid arithmetic",
			input:    `(( i = i + 1 ))`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid special vars",
			input:    `(( $# > 0 ))`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid positional",
			input:    `(( $1 > 5 ))`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid expansion",
			input:    `(( ${#arr} > 0 ))`,
			expected: []katas.Violation{},
		},
		{
			name:     "invalid simple variable",
			input:    `(( $i > 5 ))`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1073",
					Message: "Unnecessary use of `$` in arithmetic expressions. Use `(( var ))` instead of `(( $var ))`.",
					Line:    1,
					Column:  4,
				},
			},
		},
		{
			name:     "invalid multiple",
			input:    `(( $x + $y ))`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1073",
					Message: "Unnecessary use of `$` in arithmetic expressions. Use `(( var ))` instead of `(( $var ))`.",
					Line:    1,
					Column:  4,
				},
				{
					KataID:  "ZC1073",
					Message: "Unnecessary use of `$` in arithmetic expressions. Use `(( var ))` instead of `(( $var ))`.",
					Line:    1,
					Column:  9,
				},
			},
		},
		{
			name:     "valid command subst",
			input:    `(( $(date +%s) > 0 ))`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1073")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
