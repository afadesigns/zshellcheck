package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1157(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid strings with flags",
			input:    `strings -a binary`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid simple strings",
			input: `strings binary`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1157",
					Message: "Consider Zsh parameter expansion for string extraction from variables. `strings` is typically needed only for binary file analysis.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1157")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
