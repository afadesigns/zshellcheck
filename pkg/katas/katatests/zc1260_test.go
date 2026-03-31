package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1260(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid git branch -d",
			input:    `git branch -d feature`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid git branch -D",
			input: `git branch -D feature`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1260",
					Message: "Use `git branch -d` instead of `-D`. The lowercase `-d` refuses to delete unmerged branches, preventing accidental data loss.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1260")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
