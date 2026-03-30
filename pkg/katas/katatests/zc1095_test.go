package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1095(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid seq with range",
			input:    `seq 1 10`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid seq with step",
			input:    `seq 1 2 10`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid seq with single number",
			input: `seq 5`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1095",
					Message: "Use `repeat N do ... done` or `for i in {1..N}` instead of `seq N`. Zsh has built-in constructs for repetition that avoid spawning an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid seq with large number",
			input: `seq 100`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1095",
					Message: "Use `repeat N do ... done` or `for i in {1..N}` instead of `seq N`. Zsh has built-in constructs for repetition that avoid spawning an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:     "valid non-numeric argument",
			input:    `seq abc`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1095")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
