package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestCheckZC1036(t *testing.T) {
	tests := []struct {
		input    string
		expected []katas.Violation
	}{
		{
			input: `test -f file.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1036",
					Message: "Prefer `[[ ... ]]` over `test` command for conditional expressions.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1036")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
