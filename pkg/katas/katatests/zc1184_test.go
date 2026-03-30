package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1184(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid diff without -u",
			input:    `diff file1 file2`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid diff -u",
			input: `diff -u old.txt new.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1184",
					Message: "Consider `git diff` instead of `diff -u` when working in a repository. `git diff` provides better context and integration.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1184")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
