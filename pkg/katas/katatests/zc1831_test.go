// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1831(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `systemctl reload sshd`",
			input:    `systemctl reload sshd`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `systemctl status sshd`",
			input:    `systemctl status sshd`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `systemctl stop sshd`",
			input: `systemctl stop sshd`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1831",
					Message: "`systemctl stop sshd` blocks SSH — existing sessions survive but reconnects fail. `disable`/`mask` persist across reboots. Use `reload sshd` for config changes.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `systemctl mask ssh`",
			input: `systemctl mask ssh`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1831",
					Message: "`systemctl mask ssh` blocks SSH — existing sessions survive but reconnects fail. `disable`/`mask` persist across reboots. Use `reload sshd` for config changes.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1831")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
