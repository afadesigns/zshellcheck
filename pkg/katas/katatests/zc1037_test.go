package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1037(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid print",
			input:    `print -r -- "$var"`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid echo with variable",
			input: `echo $var`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1037",
					Message: "Use 'print -r --' instead of 'echo' to reliably print variable expansions.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid echo with quoted variable",
			input: `echo "$var"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1037",
					Message: "Use 'print -r --' instead of 'echo' to reliably print variable expansions.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1037")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}