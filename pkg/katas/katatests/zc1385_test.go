package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1385(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — unrelated echo",
			input:    `echo hello`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — echo $PS0",
			input: `echo $PS0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1385",
					Message: "`$PS0` is Bash-only. Zsh uses the `preexec` hook function for pre-execution prompts.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1385")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
