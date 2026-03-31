package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1196(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid grep -F",
			input:    `grep -F "hello world" file`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid grep with regex",
			input:    `grep "foo.*bar" file`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid grep with literal string",
			input: `grep "hello" file`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1196",
					Message: "Use `grep -F` for literal string matching. It skips regex compilation and is faster for fixed patterns.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1196")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
