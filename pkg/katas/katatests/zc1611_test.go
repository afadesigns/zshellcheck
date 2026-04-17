package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1611(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — Zsh ${(U)var}",
			input:    `echo "${(U)var}"`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — plain expansion",
			input:    `echo "${var}"`,
			expected: []katas.Violation{},
		},
		{
			name:  `invalid — echo "${var^^}"`,
			input: `echo "${var^^}"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1611",
					Message: "`${var^^}` / `${var,,}` — prefer Zsh `${(U)var}` / `${(L)var}` for case conversion.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  `invalid — print -r "${name,,}"`,
			input: `print -r -- "${name,,}"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1611",
					Message: "`${var^^}` / `${var,,}` — prefer Zsh `${(U)var}` / `${(L)var}` for case conversion.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1611")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
