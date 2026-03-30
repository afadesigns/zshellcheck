package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1041(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "safe printf with string literal",
			input:    `printf '%s\n' "hello"`,
			expected: []katas.Violation{},
		},
		{
			name:     "not printf command",
			input:    `echo hello`,
			expected: []katas.Violation{},
		},
		{
			name:     "printf with no args",
			input:    `printf`,
			expected: []katas.Violation{},
		},
		{
			name:  "printf with variable as format",
			input: `printf $var`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1041",
					Message: "Do not use variables in printf format string. Use `printf '..%s..' \"$var\"` instead.",
					Line:    1,
					Column:  8,
				},
			},
		},
		{
			name:     "printf with safe static format",
			input:    `printf 'hello world'`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1041")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
