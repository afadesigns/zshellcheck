package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1297(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid $0 usage",
			input:    `echo $0`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid BASH_SOURCE usage",
			input: `echo $BASH_SOURCE`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1297",
					Message: "Avoid `$BASH_SOURCE` in Zsh — use `$0` or `${(%):-%x}` instead. `BASH_SOURCE` is Bash-specific.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1297")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
