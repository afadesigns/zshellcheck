// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1416(t *testing.T) {
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
			name:  "invalid — trap 'cmd' DEBUG",
			input: `trap 'echo $BASH_COMMAND' DEBUG`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1416",
					Message: "Use Zsh `preexec() { ... }` (or `add-zsh-hook preexec`) instead of `trap 'cmd' DEBUG`. Zsh's DEBUG trap does not fire the same way as Bash's.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1416")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
