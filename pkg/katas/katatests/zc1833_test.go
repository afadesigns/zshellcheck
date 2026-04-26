// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1833(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `setopt WARN_CREATE_GLOBAL`",
			input:    `setopt WARN_CREATE_GLOBAL`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `unsetopt NOMATCH` (unrelated)",
			input:    `unsetopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `unsetopt WARN_CREATE_GLOBAL`",
			input: `unsetopt WARN_CREATE_GLOBAL`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1833",
					Message: "`unsetopt WARN_CREATE_GLOBAL` silences Zsh's warning for assignments leaking out of function scope — classic caller-variable stomping. Declare `local`/`typeset`; scope with `LOCAL_OPTIONS` if you must disable.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `setopt NO_WARN_CREATE_GLOBAL`",
			input: `setopt NO_WARN_CREATE_GLOBAL`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1833",
					Message: "`setopt NO_WARN_CREATE_GLOBAL` silences Zsh's warning for assignments leaking out of function scope — classic caller-variable stomping. Declare `local`/`typeset`; scope with `LOCAL_OPTIONS` if you must disable.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1833")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
