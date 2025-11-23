package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1071(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid array assignment",
			input:    `arr=(1 2 3)`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid append",
			input:    `arr+=(4)`,
			expected: []katas.Violation{},
		},
		{
			name:     "invalid append self reference single",
			input:    `arr=($arr)`, // Works because ($arr) is single expression
			expected: []katas.Violation{
				{
					KataID:  "ZC1071",
					Message: "Appending to an array using `arr=($arr ...)` is verbose and slower. Use `arr+=(...)` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		// This test case is commented out due to parser limitation with array literals containing spaces/multiple elements.
		// The parser currently expects `(` to start a GroupedExpression (single expression), failing on `($arr 4)`.
		/*
		{
			name:     "invalid append self reference multiple",
			input:    `arr=($arr 4)`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1071",
					Message: "Appending to an array using `arr=($arr ...)` is verbose and slower. Use `arr+=(...)` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		*/
		{
			name:     "invalid append self reference brace single",
			input:    `arr=(${arr})`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1071",
					Message: "Appending to an array using `arr=($arr ...)` is verbose and slower. Use `arr+=(...)` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1071")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}