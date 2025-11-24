package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1076(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid autoload",
			input:    `autoload -Uz my_func`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid autoload split flags",
			input:    `autoload -U -z my_func`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid autoload with other flags",
			input:    `autoload -UzX my_func`,
			expected: []katas.Violation{},
		},
		{
			name:  "missing U",
			input: `autoload -z my_func`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1076",
					Message: "Use `autoload -Uz` to ensure consistent and safe function loading.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "missing z",
			input: `autoload -U my_func`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1076",
					Message: "Use `autoload -Uz` to ensure consistent and safe function loading.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "missing flags",
			input: `autoload my_func`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1076",
					Message: "Use `autoload -Uz` to ensure consistent and safe function loading.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1076")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
