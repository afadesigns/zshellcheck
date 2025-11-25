package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1081(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid length check",
			input:    `len=${#var}`,
			expected: []katas.Violation{},
		},
		{
			name:     "invalid echo pipe wc -c",
			input:    `len=$(echo $var | wc -c)`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1081",
					Message: "Use `${#var}` to get string length. Pipeline to `wc` is inefficient.",
					Line:    1,
					Column:  17,
				},
			},
		},
		{
			name:     "invalid print pipe wc -m",
			input:    `len=$(print -r $var | wc -m)`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1081",
					Message: "Use `${#var}` to get string length. Pipeline to `wc` is inefficient.",
					Line:    1,
					Column:  21,
				},
			},
		},
		{
			name:     "wc on file (valid)",
			input:    `wc -c file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:     "cat pipe wc (valid-ish)",
			input:    `cat file | wc -c`,
			expected: []katas.Violation{}, // ZC1038 might flag cat usage, but ZC1081 shouldn't
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1081")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
