// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1622(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — Zsh flag form",
			input:    `echo "${(U)var}"`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — plain expansion",
			input:    `print -r -- "${var}"`,
			expected: []katas.Violation{},
		},
		{
			name:  `invalid — echo "${var@U}"`,
			input: `echo "${var@U}"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1622",
					Message: "`${var@U}` — prefer Zsh `${(X)var}` parameter-expansion flags (e.g. `${(U)var}` for uppercase).",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  `invalid — print "${path@Q}"`,
			input: `print -r -- "${path@Q}"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1622",
					Message: "`${var@Q}` — prefer Zsh `${(X)var}` parameter-expansion flags (e.g. `${(U)var}` for uppercase).",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1622")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
