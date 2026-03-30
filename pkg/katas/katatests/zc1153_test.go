package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1153(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid diff for viewing",
			input:    `diff file1 file2`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid diff -q for equality",
			input: `diff -q file1 file2`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1153",
					Message: "Use `cmp -s file1 file2` instead of `diff -q`. `cmp -s` is faster for equality checks as it stops at the first difference.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1153")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
