package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1120(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid pwd -P",
			input:    `pwd -P`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid bare pwd",
			input: `pwd -L`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1120",
					Message: "Use `$PWD` instead of `pwd`. Zsh maintains `$PWD` as a built-in variable, avoiding an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1120")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
