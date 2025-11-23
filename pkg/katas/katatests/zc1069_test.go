package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1069(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid local in function",
			input:    `my_func() { local x=1; }`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid typeset global",
			input:    `typeset x=1`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid local global",
			input: `local x=1`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1069",
					Message: "`local` can only be used inside functions. " +
						"Use `typeset`, `declare`, or just assignment for global variables.",
					Line:   1,
					Column: 1,
				},
			},
		},
		{
			name:  "invalid local in if block (global)",
			input: `if true; then local x=1; fi`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1069",
					Message: "`local` can only be used inside functions. " +
						"Use `typeset`, `declare`, or just assignment for global variables.",
					Line:   1,
					Column: 15,
				},
			},
		},
		{
			name:     "valid local in nested function",
			input:    `outer() { inner() { local x=1; }; }`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid local in subshell (global)",
			input: `( local x=1 )`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1069",
					Message: "`local` can only be used inside functions. " +
						"Use `typeset`, `declare`, or just assignment for global variables.",
					Line:   1,
					Column: 3,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1069")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
