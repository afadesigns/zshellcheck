package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1517(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — print -P with literal",
			input:    `print -P "%F{red}hello%f"`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — print without -P",
			input:    `print "$var"`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — print -P with single-quoted var (no interpolation)",
			input:    `print -P '$var'`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — print -P \"$var\"",
			input: `print -P "$var"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1517",
					Message: "`print -P \"$var\"` expands prompt escapes inside the variable — use `${(V)var}` / `${(q-)var}` to neutralize metacharacters, or drop -P.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — print -P $msg (unquoted)",
			input: `print -P $msg`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1517",
					Message: "`print -P $msg` expands prompt escapes inside the variable — use `${(V)var}` / `${(q-)var}` to neutralize metacharacters, or drop -P.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1517")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
