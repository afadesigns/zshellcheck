package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestCheckZC1034(t *testing.T) {
	tests := []struct {
		input    string
		expected []katas.Violation
	}{
		{
			input: `which ls`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1034",
					Message: "Use `command -v` instead of `which` for portability.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			input:    `command -v ls`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1034")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
