// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1311(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid compdef usage",
			input:    `compdef _git git`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid complete usage",
			input: `complete -F _git git`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1311",
					Message: "Avoid `complete` in Zsh — it is a Bash builtin. Use `compdef` for Zsh completion registration.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1311")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
