package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1044(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "cd with error handling",
			input:    `cd /tmp || exit 1`,
			expected: []katas.Violation{},
		},
		{
			name:  "unchecked cd",
			input: `cd /tmp`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1044",
					Message: "Use `cd ... || return` (or `exit`) in case cd fails.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:     "no cd command",
			input:    `echo hello`,
			expected: []katas.Violation{},
		},
		{
			name:     "cd in if condition",
			input:    `if cd /tmp; then echo ok; fi`,
			expected: []katas.Violation{},
		},
		{
			name:  "cd with && chain",
			input: `cd /tmp && echo ok`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1044",
					Message: "Use `cd ... || return` (or `exit`) in case cd fails.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:     "cd negated in condition",
			input:    `if ! cd /tmp; then echo fail; fi`,
			expected: []katas.Violation{},
		},
		{
			name:  "cd in while loop body",
			input: `while true; do cd /tmp; done`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1044",
					Message: "Use `cd ... || return` (or `exit`) in case cd fails.",
					Line:    1,
					Column:  16,
				},
			},
		},
		{
			name:  "cd in for loop body",
			input: `for d in /tmp /var; do cd $d; done`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1044",
					Message: "Use `cd ... || return` (or `exit`) in case cd fails.",
					Line:    1,
					Column:  24,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1044")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
