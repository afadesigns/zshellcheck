package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1191(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid command run",
			input:    `command ls`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid command -v check",
			input: `command -v git`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1191",
					Message: "Use `(( $+commands[cmd] ))` instead of `command -v cmd`. Zsh `$commands` array provides instant command lookups.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1191")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
