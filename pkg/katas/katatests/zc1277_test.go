package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1277(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid tr for other transformations",
			input:    `tr -d '\n'`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid Zsh :l modifier",
			input:    `echo ${var:l}`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid tr uppercase to lowercase POSIX class",
			input: `tr '[:upper:]' '[:lower:]'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1277",
					Message: "Use Zsh parameter expansion `${var:l}` for lowercase conversion instead of `tr`. The `:l` modifier is built-in.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid tr lowercase to uppercase range",
			input: `tr '[a-z]' '[A-Z]'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1277",
					Message: "Use Zsh parameter expansion `${var:u}` for uppercase conversion instead of `tr`. The `:u` modifier is built-in.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1277")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
