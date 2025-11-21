package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1003(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid arithmetic test",
			input:    `(( 1 > 0 ))`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid arithmetic test",
			input: `[ 1 -gt 0 ]`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1003",
					Message: "Prefer [[ over [ for tests. " +
						"[[ is a Zsh keyword that offers safer and more powerful conditional expressions.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1003")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
