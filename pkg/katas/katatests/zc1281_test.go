package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1281(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid sort -u usage",
			input:    `sort -u file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid uniq with flags",
			input:    `uniq -c file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid uniq with file",
			input: `uniq file.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1281",
					Message: "Use `sort -u` instead of `sort | uniq`. The `-u` flag deduplicates in a single pass.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1281")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
