package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1068(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid hook registration",
			input:    `autoload -Uz add-zsh-hook; add-zsh-hook precmd my_precmd`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid normal function",
			input:    `my_func() { echo hello; }`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid precmd definition",
			input: `precmd() { echo "prompt"; }`,
			expected: []katas.Violation{
				{
					KataID: "ZC1068",
					Message: "Defining `precmd` directly overwrites existing hooks. " +
						"Use `autoload -Uz add-zsh-hook; add-zsh-hook precmd my_func` instead.",
					Line:   1,
					Column: 1,
				},
			},
		},
		{
			name:  "invalid chpwd definition",
			input: `function chpwd() { ls; }`,
			expected: []katas.Violation{
				{
					KataID: "ZC1068",
					Message: "Defining `chpwd` directly overwrites existing hooks. " +
						"Use `autoload -Uz add-zsh-hook; add-zsh-hook chpwd my_func` instead.",
					Line:   1,
					// Start of "function" keyword usually, or name depending on parser.
					Column: 1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1068")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
