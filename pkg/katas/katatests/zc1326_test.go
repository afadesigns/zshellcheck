package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1326(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid HISTFILE usage",
			input:    `echo $HISTFILE`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid HISTTIMEFORMAT usage",
			input: `echo $HISTTIMEFORMAT`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1326",
					Message: "Avoid `$HISTTIMEFORMAT` in Zsh — use `setopt EXTENDED_HISTORY` and `fc -li` instead.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1326")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
