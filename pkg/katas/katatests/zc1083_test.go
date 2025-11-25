package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1083(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid brace expansion",
			input:    `echo {1..10}`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid brace expansion with list",
			input:    `echo {a,b,c}`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid variable after braces",
			input:    `echo {1..10}$var`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid variable before braces",
			input:    `echo $var{1..10}`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid variable inside list expansion",
			input:    `echo {a,b,$var}`,
			expected: []katas.Violation{},
		},
		{
			name:     "invalid variable as range end",
			input:    `echo {1..$n}`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1083",
					Message: "Brace expansion limits cannot be variables. `{...$var...}` is treated as a literal string. Use `seq` or `for ((...))`.",
					Line:    1,
					Column:  6,
				},
			},
		},
		{
			name:     "invalid variable as range start",
			input:    `echo {$n..10}`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1083",
					Message: "Brace expansion limits cannot be variables. `{...$var...}` is treated as a literal string. Use `seq` or `for ((...))`.",
					Line:    1,
					Column:  6,
				},
			},
		},
		{
			name:     "invalid variable as range start and end",
			input:    `echo {$min..$max}`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1083",
					Message: "Brace expansion limits cannot be variables. `{...$var...}` is treated as a literal string. Use `seq` or `for ((...))`.",
					Line:    1,
					Column:  6,
				},
			},
		},
		{
			name:     "invalid command substitution in range",
			input:    `echo {1..$(echo 10)}`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1083",
					Message: "Brace expansion limits cannot be variables. `{...$var...}` is treated as a literal string. Use `seq` or `for ((...))`.",
					Line:    1,
					Column:  6,
				},
			},
		},
		{
			name:     "invalid quoted brace expansion with variable",
			input:    `echo "{1..$n}"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1083",
					Message: "Brace expansion limits cannot be variables. `{...$var...}` is treated as a literal string. Use `seq` or `for ((...))`.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1083")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
