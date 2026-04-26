// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1417(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — trap 'cleanup' EXIT",
			input:    `trap 'cleanup' EXIT`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — trap 'cmd' RETURN",
			input: `trap 'print done' RETURN`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1417",
					Message: "Prefer Zsh `TRAPRETURN() { ... }` function over `trap 'cmd' RETURN`. Named-function form is more idiomatic in Zsh.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1417")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
