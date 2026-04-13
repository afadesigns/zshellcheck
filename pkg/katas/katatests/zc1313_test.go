package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1313(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid aliases usage",
			input:    `echo $aliases`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid BASH_ALIASES usage",
			input: `echo $BASH_ALIASES`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1313",
					Message: "Avoid `$BASH_ALIASES` in Zsh — use the `aliases` associative array instead.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1313")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
