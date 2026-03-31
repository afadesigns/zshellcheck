package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1259(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid grep -rI",
			input:    `grep -rI pattern .`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid grep without -r",
			input:    `grep pattern file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid grep -r without -I",
			input: `grep -r pattern .`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1259",
					Message: "Use `grep -I` with recursive search to skip binary files. Without `-I`, grep may produce garbled binary output.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1259")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
