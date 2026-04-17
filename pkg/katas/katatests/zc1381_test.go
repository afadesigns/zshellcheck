package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1381(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — echo $words (Zsh compsys)",
			input:    `echo $words`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — echo $COMP_WORDS",
			input: `echo $COMP_WORDS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1381",
					Message: "Bash `$COMP_*` completion variables do not exist in Zsh. Use `$words` (array of tokens), `$CURRENT` (cursor index), `$BUFFER`, or the `_arguments`/`_values` helpers from `compsys`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — echo $COMP_CWORD",
			input: `echo $COMP_CWORD`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1381",
					Message: "Bash `$COMP_*` completion variables do not exist in Zsh. Use `$words` (array of tokens), `$CURRENT` (cursor index), `$BUFFER`, or the `_arguments`/`_values` helpers from `compsys`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1381")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
