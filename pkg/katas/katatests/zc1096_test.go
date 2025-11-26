package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1096(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:  "bc usage",
			input: `echo "1.5 + 2.5" | bc`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1096",
					Message: "Zsh supports floating point arithmetic natively. You often don't need `bc`.",
					Line:    1,
					Column:  20,
				},
			},
		},
		{
			name:  "valid arithmetic",
			input: `(( 1.5 + 2.5 ))`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1096")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
