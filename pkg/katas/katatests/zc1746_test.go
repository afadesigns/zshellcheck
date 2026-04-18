package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1746(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `sysctl -w kernel.randomize_va_space=2` (default ASLR)",
			input:    `sysctl -w kernel.randomize_va_space=2`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `sysctl kernel.randomize_va_space`",
			input:    `sysctl kernel.randomize_va_space`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `sysctl -w kernel.randomize_va_space=0`",
			input: `sysctl -w kernel.randomize_va_space=0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1746",
					Message: "`sysctl kernel.randomize_va_space=0` weakens ASLR — absolute-address exploits become deterministic (stack overflows, ROP). Keep `kernel.randomize_va_space=2` outside a sandboxed debug context.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `sysctl kernel.randomize_va_space=1`",
			input: `sysctl kernel.randomize_va_space=1`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1746",
					Message: "`sysctl kernel.randomize_va_space=1` weakens ASLR — absolute-address exploits become deterministic (stack overflows, ROP). Keep `kernel.randomize_va_space=2` outside a sandboxed debug context.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1746")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
