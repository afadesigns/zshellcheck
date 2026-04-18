package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1673(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — stty echo (restore)",
			input:    `stty echo`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — stty raw",
			input:    `stty raw`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — stty -echo",
			input: `stty -echo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1673",
					Message: "`stty -echo` to mask password entry is fragile — a crash leaves the terminal echo-off. Use `read -s VAR` (Zsh / Bash 4+) instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1673")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
