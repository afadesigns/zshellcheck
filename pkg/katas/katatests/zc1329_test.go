package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1329(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid HISTSIZE usage",
			input:    `echo $HISTSIZE`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid HISTIGNORE usage",
			input: `echo $HISTIGNORE`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1329",
					Message: "Avoid `$HISTIGNORE` in Zsh — use `zshaddhistory` hook for history filtering instead.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1329")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
