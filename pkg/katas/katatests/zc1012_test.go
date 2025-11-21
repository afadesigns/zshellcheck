package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestCheckZC1012(t *testing.T) {
	tests := []struct {
		input    string
		expected []katas.Violation
	}{
		{
			input: "`ls`",
			expected: []katas.Violation{
				{
					KataID:  "ZC1012",
					Message: "Use `$(command)` instead of backticks for command substitution.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			input:    `$(ls)`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1012")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}