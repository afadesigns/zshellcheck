package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1126(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid sort -u",
			input:    `sort -u file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid sort | uniq -c",
			input:    `sort file | uniq -c`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid sort | uniq",
			input: `sort file | uniq`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1126",
					Message: "Use `sort -u` instead of `sort | uniq`. Combining into one command avoids an unnecessary pipeline.",
					Line:    1,
					Column:  11,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1126")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
