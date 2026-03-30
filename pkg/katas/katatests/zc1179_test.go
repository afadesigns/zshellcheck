package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1179(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid date without format",
			input:    `date -u`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid date with format",
			input: `date '+%Y-%m-%d'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1179",
					Message: "Use `strftime` (via `zmodload zsh/datetime`) instead of `date +%Y-%m-%d`. Zsh date formatting avoids spawning an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1179")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
