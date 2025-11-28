package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1106(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:  "set -x usage",
			input: `set -x`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1106",
					Message: "Avoid `set -x` in production scripts to prevent sensitive data exposure.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "set -eux usage",
			input: `set -eux`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1106",
					Message: "Avoid `set -x` in production scripts to prevent sensitive data exposure.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:     "valid set usage",
			input:    `set +x`,
			expected: []katas.Violation{},
		},
		{
			name:     "other set options",
			input:    `set -e`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1106")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
