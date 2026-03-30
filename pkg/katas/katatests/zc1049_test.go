package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1049(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:  "alias usage",
			input: `alias ll='ls -la'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1049",
					Message: "Prefer functions over aliases. Aliases are expanded at parse time and can behave unexpectedly in scripts.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:     "function instead of alias",
			input:    `echo hello`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1049")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
