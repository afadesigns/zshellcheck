package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1380(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — echo $HISTORY_IGNORE (Zsh)",
			input:    `echo $HISTORY_IGNORE`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — echo $HISTIGNORE",
			input: `echo $HISTIGNORE`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1380",
					Message: "`$HISTIGNORE` is Bash-only. In Zsh use `$HISTORY_IGNORE` (underscored) for the same history-pattern filter.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1380")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
