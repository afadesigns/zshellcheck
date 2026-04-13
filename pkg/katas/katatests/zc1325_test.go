package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1325(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid PS1 usage",
			input:    `echo $PS1`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid PS0 usage",
			input: `echo $PS0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1325",
					Message: "Avoid `$PS0` in Zsh — use the `preexec` hook function instead. `PS0` is Bash 4.4+ specific.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1325")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
