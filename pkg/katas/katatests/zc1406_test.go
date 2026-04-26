// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1406(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — xargs without -P",
			input:    `xargs -n 1 echo`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — xargs -P 4",
			input: `xargs -P 4 cmd`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1406",
					Message: "Consider `zargs -P N` (autoload -Uz zargs) instead of `xargs -P N`. Parallel execution with Zsh functions in scope — no subshell-per-item.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — xargs -P4 attached",
			input: `xargs -P4 cmd`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1406",
					Message: "Consider `zargs -P N` (autoload -Uz zargs) instead of `xargs -P N`. Parallel execution with Zsh functions in scope — no subshell-per-item.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1406")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
