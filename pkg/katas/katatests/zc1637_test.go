// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1637(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — unrelated command",
			input:    `export FOO=bar`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — readonly FOO=bar",
			input: `readonly FOO=bar`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1637",
					Message: "`readonly` works but `typeset -r NAME=value` is the Zsh-native form and composes with other typeset flags.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — readonly MAX_RETRIES=5",
			input: `readonly MAX_RETRIES=5`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1637",
					Message: "`readonly` works but `typeset -r NAME=value` is the Zsh-native form and composes with other typeset flags.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1637")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
