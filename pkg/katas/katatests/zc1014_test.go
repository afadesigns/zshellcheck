package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestCheckZC1014(t *testing.T) {
	tests := []struct {
		input    string
		expected []katas.Violation
	}{
		{
			input: `git checkout my-branch`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1014",
					Message: "Use `git switch` or `git restore` instead of the ambiguous `git checkout`.",
										Line:    1,
										Column:  1,
									},
								},
							},
							{
								input:    `[ -f file ]`,
			expected: []katas.Violation{},
		},
		{
			input:    `git restore my-file`,
			expected: []katas.Violation{},
		},
		{
			input:    `git commit -m "message"`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1014")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
