// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1774(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `setopt EXTENDED_GLOB` (unrelated option)",
			input:    `setopt EXTENDED_GLOB`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `set -e` (unrelated short option)",
			input:    `set -e`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt GLOB_SUBST`",
			input: `setopt GLOB_SUBST`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1774",
					Message: "`setopt GLOB_SUBST` enables `GLOB_SUBST` — every unquoted `$var` expansion is rescanned as a glob pattern. User-controlled data becomes a filesystem query. Scope this in a subshell / function, or use explicit expansion flags.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `setopt globsubst` (Zsh lower/underscore folded)",
			input: `setopt globsubst`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1774",
					Message: "`setopt globsubst` enables `GLOB_SUBST` — every unquoted `$var` expansion is rescanned as a glob pattern. User-controlled data becomes a filesystem query. Scope this in a subshell / function, or use explicit expansion flags.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `set -o GLOB_SUBST`",
			input: `set -o GLOB_SUBST`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1774",
					Message: "`set -o GLOB_SUBST` enables `GLOB_SUBST` — every unquoted `$var` expansion is rescanned as a glob pattern. User-controlled data becomes a filesystem query. Scope this in a subshell / function, or use explicit expansion flags.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1774")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
