package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestCheckZC1012(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:  "read without flags",
			input: `read line`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1012",
					Message: "Use `read -r` to read input without interpreting backslashes.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:     "read with -r",
			input:    `read -r line`,
			expected: []katas.Violation{},
		},
        {
			name:     "read with -er",
			input:    `read -er line`, // heuristic support
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1012")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
