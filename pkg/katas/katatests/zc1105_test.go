package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1105(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:  "nested arithmetic expansion",
			input: `(( result = $((1+1)) + 2 ))`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1105",
					Message: "Avoid nested arithmetic expansions. Use intermediate variables for clarity.",
					Line:    1,
					Column:  2,
				},
			},
		},
		{
			name:     "simple arithmetic expansion",
			input:    `result=$((1+1))`,
			expected: []katas.Violation{},
		},
		{
			name:     "double parenthesis arithmetic command",
			input:    `(( 1 + 1 ))`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1105")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
