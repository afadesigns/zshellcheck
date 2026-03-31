package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1237(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid git clean -n",
			input:    `git clean -n`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid git clean -fd",
			input: `git clean -fd`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1237",
					Message: "Use `git clean -n` first to preview removals before `git clean -fd`. Forced clean permanently deletes untracked files.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1237")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
