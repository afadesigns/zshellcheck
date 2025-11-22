package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestCheckZC1030(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			            name:  "echo with a simple string",
			            input: `echo "hello"`,
			            expected: []katas.Violation{				{
					KataID:  "ZC1030",
					Message: "Use `printf` for more reliable and portable string formatting instead of `echo`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:     "printf with a simple string",
			input:    `printf "hello"`,
			expected: []katas.Violation{},
		},
		{
			name:     "echo with a variable",
			input:    `echo "$foo"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1030",
					Message: "Use `printf` for more reliable and portable string formatting instead of `echo`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1030")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}