package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1321(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid variable",
			input:    `echo $MY_FD`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid BASH_XTRACEFD usage",
			input: `echo $BASH_XTRACEFD`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1321",
					Message: "Avoid `$BASH_XTRACEFD` in Zsh — it is undefined. Redirect stderr directly for xtrace output.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1321")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
