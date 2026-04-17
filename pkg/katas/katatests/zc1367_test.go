package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1367(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — printf normal format",
			input:    `printf '%s\n' hello`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — printf %(fmt)T",
			input: `printf '%(%Y-%m-%d)T\n' 1700000000`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1367",
					Message: "Use Zsh `strftime fmt seconds` (from `zsh/datetime`) instead of Bash `printf '%(fmt)T' seconds`. Same formatting, more readable, no Bash-version gating.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1367")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
