package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1148(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid compdef",
			input:    `compdef _git git`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid compctl",
			input: `compctl -K _my_func mycommand`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1148",
					Message: "Use `compdef` instead of `compctl`. The `compctl` system is deprecated; use `compinit` and `compdef` for modern Zsh completions.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1148")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
