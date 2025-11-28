package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1097(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid local loop variable",
			input:    `my_func() { local i; for i in 1 2; do echo $i; done; }`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid global loop variable",
			input: `my_func() { for i in 1 2; do echo $i; done; }`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1097",
					Message: "Loop variable 'i' is used without 'local'. It will be global. Use `local i` before the loop.",
					Line:    1,
					Column:  17,
				},
			},
		},
		{
			name:  "invalid global loop variable (implicit in)",
			input: `my_func() { for i; do echo $i; done; }`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1097",
					Message: "Loop variable 'i' is used without 'local'. It will be global. Use `local i` before the loop.",
					Line:    1,
					Column:  17,
				},
			},
		},
		{
			name:     "valid typeset loop variable",
			input:    `my_func() { typeset i; for i in 1 2; do echo $i; done; }`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid integer loop variable",
			input:    `my_func() { integer i; for i in 1 2; do echo $i; done; }`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid local loop variable in nested block",
			input:    `my_func() { if true; then local i; for i in 1 2; do echo $i; done; fi; }`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid local inside loop (too late)",
			input: `my_func() { for i in 1 2; do local i; echo $i; done; }`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1097",
					Message: "Loop variable 'i' is used without 'local'. It will be global. Use `local i` before the loop.",
					Line:    1,
					Column:  17,
				},
			},
		},
		{
			name:     "valid arithmetic for loop (C-style)",
			input:    `my_func() { for ((i=0; i<10; i++)); do echo $i; done; }`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1097")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
