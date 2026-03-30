package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1119(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid date with format",
			input:    `date "+%Y-%m-%d"`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid date no args",
			input:    `date -u`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid date +%s",
			input: `date '+%s'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1119",
					Message: "Use `$EPOCHSECONDS` or `$EPOCHREALTIME` (via `zmodload zsh/datetime`) instead of `date +%s`. Avoids spawning an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1119")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
