// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1979(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt HIST_FCNTL_LOCK` (keeps default off)",
			input:    `unsetopt HIST_FCNTL_LOCK`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NO_HIST_FCNTL_LOCK`",
			input:    `setopt NO_HIST_FCNTL_LOCK`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt HIST_FCNTL_LOCK`",
			input: `setopt HIST_FCNTL_LOCK`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1979",
					Message: "`setopt HIST_FCNTL_LOCK` routes `$HISTFILE` locking through POSIX `fcntl()` — on NFS home directories a hung `rpc.lockd` freezes every other shell at the next prompt. Keep off; enable only when `$HISTFILE` is on a local fs.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_HIST_FCNTL_LOCK`",
			input: `unsetopt NO_HIST_FCNTL_LOCK`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1979",
					Message: "`unsetopt NO_HIST_FCNTL_LOCK` routes `$HISTFILE` locking through POSIX `fcntl()` — on NFS home directories a hung `rpc.lockd` freezes every other shell at the next prompt. Keep off; enable only when `$HISTFILE` is on a local fs.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1979")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
