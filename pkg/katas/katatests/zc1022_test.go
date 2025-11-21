package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestCheckZC1022(t *testing.T) {
	tests := []struct {
		input    string
		expected []katas.Violation
	}{
		{
			input: `let x=1+1`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1022",
					Message: "Use `$((...))` for arithmetic expansion instead of `let`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			input:    `x=$((1+1))`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1022")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
