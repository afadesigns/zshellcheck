package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1389(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — echo $hosts (Zsh)",
			input:    `echo $hosts`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — echo $HOSTFILE",
			input: `echo $HOSTFILE`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1389",
					Message: "`$HOSTFILE` is Bash-only. Zsh reads hostnames for completion from the `$hosts` array (lowercase).",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1389")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
