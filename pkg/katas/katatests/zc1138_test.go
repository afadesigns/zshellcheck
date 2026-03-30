package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1138(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid printf with format",
			input:    `printf "Hello %s" name`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid printf newline pattern",
			input: `printf '%s\n' item1`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1138",
					Message: "Use `print -l` instead of `printf '%s\\n'` for printing elements one per line. `print -l` is a Zsh builtin optimized for this pattern.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1138")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
