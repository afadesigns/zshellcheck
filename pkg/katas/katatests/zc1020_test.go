package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestCheckZC1020(t *testing.T) {
	tests := []struct {
		input    string
		expected []katas.Violation
	}{
		{
			input: `test 1 -eq 1`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1020",
					Message: "Use `[[ ... ]]` for tests instead of `test`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			input:    `[[ 1 -eq 1 ]]`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1020")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
