// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1318(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid commands usage",
			input:    `echo $commands`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid BASH_CMDS usage",
			input: `echo $BASH_CMDS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1318",
					Message: "Avoid `$BASH_CMDS` in Zsh — use the `$commands` hash for command path lookups instead.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1318")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
