package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestCheckZC1032(t *testing.T) {
	tests := []struct {
		input    string
		expected []katas.Violation
	}{
		{
			input: `let i=i+1`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1032",
					Message: "Use `(( i++ ))` for C-style incrementing instead of `let i=i+1`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			input:    `(( i++ ))`,
			expected: []katas.Violation{},
		},
		{
			input:    `let i=j+1`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1032")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
