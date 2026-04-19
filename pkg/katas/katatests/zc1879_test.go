package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1879(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `setopt BAD_PATTERN` (explicit default)",
			input:    `setopt BAD_PATTERN`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `unsetopt NOMATCH` (unrelated)",
			input:    `unsetopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `unsetopt BAD_PATTERN`",
			input: `unsetopt BAD_PATTERN`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1879",
					Message: "`unsetopt BAD_PATTERN` silences `bad pattern` errors — `rm [abc` tries to remove a literal `[abc`, broken `case` arms stop firing. Keep the option on; quote one-off patterns or scope with `LOCAL_OPTIONS`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `setopt NO_BAD_PATTERN`",
			input: `setopt NO_BAD_PATTERN`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1879",
					Message: "`setopt NO_BAD_PATTERN` silences `bad pattern` errors — `rm [abc` tries to remove a literal `[abc`, broken `case` arms stop firing. Keep the option on; quote one-off patterns or scope with `LOCAL_OPTIONS`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1879")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
