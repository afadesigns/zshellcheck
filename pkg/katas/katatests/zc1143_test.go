package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1143(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid setopt",
			input:    `set -u`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid set -e",
			input: `set -e`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1143",
					Message: "Avoid `set -e`. It has surprising behavior with conditionals and subshells in Zsh. Use explicit error handling with `cmd || return 1` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1143")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
