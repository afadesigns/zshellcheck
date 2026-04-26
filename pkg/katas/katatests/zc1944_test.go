// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1944(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt IGNORE_EOF` (explicit default)",
			input:    `unsetopt IGNORE_EOF`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt EMACS` (unrelated)",
			input:    `setopt EMACS`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt IGNORE_EOF`",
			input: `setopt IGNORE_EOF`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1944",
					Message: "`setopt IGNORE_EOF` makes `Ctrl-D` stop terminating the shell — subshells, sudo holds, SSH tunnels linger after the parent left. Keep off; use `TMOUT=NN` for a timed stale-tty exit if needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_IGNORE_EOF`",
			input: `unsetopt NO_IGNORE_EOF`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1944",
					Message: "`unsetopt NO_IGNORE_EOF` makes `Ctrl-D` stop terminating the shell — subshells, sudo holds, SSH tunnels linger after the parent left. Keep off; use `TMOUT=NN` for a timed stale-tty exit if needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1944")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
