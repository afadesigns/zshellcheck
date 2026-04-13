package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1336(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid typeset -x usage",
			input:    `typeset -x PATH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid printenv usage",
			input: `printenv HOME`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1336",
					Message: "Avoid `printenv` in Zsh — use `typeset -x` or `export` to list environment variables.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1336")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
