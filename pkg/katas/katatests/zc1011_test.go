package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestCheckZC1011(t *testing.T) {
	tests := []struct {
		input    string
		expected []katas.Violation
	}{
		{
			input: `git rev-parse HEAD`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1011",
					Message: "Avoid using `git` plumbing commands in scripts. They are not guaranteed to be stable.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			input:    `git branch`,
			expected: []katas.Violation{},
		},
		{
			input:    `ls -l`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1011")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
