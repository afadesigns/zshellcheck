package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestCheckZC1013(t *testing.T) {
	tests := []struct {
		input    string
		expected []katas.Violation
	}{
		{
			input: `let x = 1`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1013",
					Message: "Use `((...))` for arithmetic operations instead of `let`.",
					 Line:    1,
					 Column:  1,
					  },
					  },
					 },
					 {
					  input:    `[[ -f file ]]`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1013")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
