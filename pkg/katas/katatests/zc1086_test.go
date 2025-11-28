package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1086(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid function definition",
			input:    `my_func() { echo hello; }`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid function definition with keyword",
			input: `function my_func { echo hello; }`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1086",
					Message: "Prefer `func() { ... }` over `function func { ... }` for portability and consistency.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid function definition with keyword and parens",
			input: `function my_func() { echo hello; }`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1086",
					Message: "Prefer `func() { ... }` over `function func { ... }` for portability and consistency.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1086")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
