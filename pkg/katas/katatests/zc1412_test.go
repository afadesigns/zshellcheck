package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1412(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — echo $candidates",
			input:    `echo $candidates`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — echo $COMPREPLY",
			input: `echo $COMPREPLY`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1412",
					Message: "`$COMPREPLY` is a Bash-only completion output array. In Zsh compsys use `compadd -- candidate1 candidate2`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1412")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
