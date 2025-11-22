package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestCheckZC1017(t *testing.T) {
	tests := []struct {
		input    string
		expected []katas.Violation
	}{
		{
			input: `print "hello"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1017",
					Message: "Use `print -r` to print strings literally.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			input:    `print -r "hello"`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1017")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
