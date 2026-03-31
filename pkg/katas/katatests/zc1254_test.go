package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1254(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid git commit",
			input:    `git commit -m "feat: add feature"`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid git commit --amend",
			input: `git commit --amend -m "fix"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1254",
					Message: "Avoid `git commit --amend` on shared branches — it rewrites history. Use `git commit --fixup` or create a new commit instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1254")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
