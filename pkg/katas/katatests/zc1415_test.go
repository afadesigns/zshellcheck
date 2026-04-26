// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1415(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — trap 'cmd' EXIT",
			input:    `trap 'cleanup' EXIT`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — trap 'cmd' ERR",
			input: `trap 'echo oops' ERR`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1415",
					Message: "Prefer Zsh `TRAPZERR() { ... }` function over `trap 'cmd' ERR`. The named-function form is more idiomatic and composable in Zsh.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1415")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
