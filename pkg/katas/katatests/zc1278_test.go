package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1278(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid arithmetic expansion",
			input:    `echo $(( 1 + 2 ))`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid other command",
			input:    `echo hello`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid expr usage",
			input: `expr 1 + 2`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1278",
					Message: "Use Zsh arithmetic expansion `$(( ))` instead of `expr`. It is built-in and avoids forking.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1278")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
