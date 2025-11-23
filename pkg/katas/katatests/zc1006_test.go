package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestCheckZC1006(t *testing.T) {
	tests := []struct {
		input    string
		expected []katas.Violation
	}{
		{
			input: `test 1 -eq 1`,
			expected: []katas.Violation{
				{
					KataID: "ZC1006",
					Message: "Prefer [[ over test for tests. " +
						"[[ is a Zsh keyword that offers safer and more powerful conditional expressions.",
					Line:   1,
					Column: 1,
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
			violations := testutil.Check(tt.input, "ZC1006")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
