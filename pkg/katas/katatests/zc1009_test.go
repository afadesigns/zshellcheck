package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestCheckZC1009(t *testing.T) {
	tests := []struct {
		input    string
		expected []katas.Violation
	}{
		{
			input: `expr 1 + 1`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1009",
					Message: "Use `((...))` for C-style arithmetic instead of `expr`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			input:    `echo "hello"`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1009")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
