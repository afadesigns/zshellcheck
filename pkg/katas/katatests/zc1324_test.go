package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1324(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid precmd usage",
			input:    `echo $precmd`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid PROMPT_COMMAND usage",
			input: `echo $PROMPT_COMMAND`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1324",
					Message: "Avoid `$PROMPT_COMMAND` in Zsh — use the `precmd` hook function instead. `PROMPT_COMMAND` is Bash-specific.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1324")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
