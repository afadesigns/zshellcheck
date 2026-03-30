package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1118(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid echo without -n",
			input:    `echo hello`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid print -rn",
			input:    `print -rn hello`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid echo -n",
			input: `echo -n hello`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1118",
					Message: "Use `print -rn` instead of `echo -n`. `echo -n` behavior varies across shells; `print -rn` is the reliable Zsh idiom.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1118")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
