// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1655(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — read -k 1 (Zsh)",
			input:    `read -k 1 var`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — read -r line",
			input:    `read -r line`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — read -n 1 char",
			input: `read -n 1 char`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1655",
					Message: "`read -n N` is Bash syntax for \"read N characters\". Zsh's `-n` means \"drop trailing newline\" with no count. Use `read -k N var` in Zsh.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — read -n5 var",
			input: `read -n5 var`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1655",
					Message: "`read -n N` is Bash syntax for \"read N characters\". Zsh's `-n` means \"drop trailing newline\" with no count. Use `read -k N var` in Zsh.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1655")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
