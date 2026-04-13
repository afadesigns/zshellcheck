package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1312(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid compadd usage",
			input:    `compadd foo bar baz`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid compgen usage",
			input: `compgen -W "foo bar" -`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1312",
					Message: "Avoid `compgen` in Zsh — it is a Bash builtin. Use `compadd` or Zsh completion functions instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1312")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
