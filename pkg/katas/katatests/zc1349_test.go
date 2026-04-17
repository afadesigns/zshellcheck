package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1349(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — expr arithmetic (not length)",
			input:    `expr 1 + 2`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — expr length with var",
			input: `expr length "$s"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1349",
					Message: "Use `${#var}` instead of `expr length \"$var\"` for string length. Parameter expansion avoids spawning an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — expr length with literal",
			input: `expr length hello`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1349",
					Message: "Use `${#var}` instead of `expr length \"$var\"` for string length. Parameter expansion avoids spawning an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1349")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
