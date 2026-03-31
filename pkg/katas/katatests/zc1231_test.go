package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1231(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid git clone --depth 1",
			input:    `git clone --depth 1 https://github.com/user/repo`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid git clone full",
			input: `git clone https://github.com/user/repo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1231",
					Message: "Consider `git clone --depth 1` in scripts. Full clones download entire history which is unnecessary for builds and CI.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1231")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
