package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1002(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid command substitution",
			input:    `x=$(ls)`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid backticks",
			input: `x=` + "`ls`",
			expected: []katas.Violation{
				{
					KataID: "ZC1002",
					Message: "Use $(...) instead of backticks for command substitution. " +
						"The `$(...)` syntax is more readable and can be nested easily.",
					Line:   1,
					Column: 3,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1002")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
