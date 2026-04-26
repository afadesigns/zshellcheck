// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1180(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid pgrep with -u flag",
			input:    `pgrep -u root sshd`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid simple pgrep",
			input: `pgrep myprocess`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1180",
					Message: "For own background jobs, use Zsh job control (`jobs`, `kill %N`) instead of `pgrep`. Job control is more precise for script-spawned processes.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid pkill",
			input: `pkill -f myserver`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1180",
					Message: "For own background jobs, use Zsh job control (`jobs`, `kill %N`) instead of `pkill`. Job control is more precise for script-spawned processes.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1180")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
