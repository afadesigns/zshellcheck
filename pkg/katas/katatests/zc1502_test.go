package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1502(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — grep literal pattern",
			input:    `grep "hello" file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — grep -- \"$var\" (end-of-flags marker)",
			input:    `grep -- "$pattern" file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — grep \"$var\" file (no --)",
			input: `grep "$pattern" file.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1502",
					Message: "Variable `\"$pattern\"` used as pattern without `--` end-of-flags marker — attacker-controlled leading `-` becomes a flag. Write `grep -- \"$var\"`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — rg $pattern files (unquoted)",
			input: `rg $pattern files`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1502",
					Message: "Variable `$pattern` used as pattern without `--` end-of-flags marker — attacker-controlled leading `-` becomes a flag. Write `grep -- \"$var\"`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1502")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
