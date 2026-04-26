// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1778(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `systemctl start foo.service`",
			input:    `systemctl start foo.service`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `systemctl link /etc/systemd/system/foo.service` (immutable path)",
			input:    `systemctl link /etc/systemd/system/foo.service`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `systemctl link /tmp/foo.service`",
			input: `systemctl link /tmp/foo.service`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1778",
					Message: "`systemctl link /tmp/foo.service` keeps the unit in a mutable path — a tamper later changes what systemd runs. Copy the unit into `/etc/systemd/system/` with root-only perms.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `systemctl link /home/user/build/foo.service`",
			input: `systemctl link /home/user/build/foo.service`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1778",
					Message: "`systemctl link /home/user/build/foo.service` keeps the unit in a mutable path — a tamper later changes what systemd runs. Copy the unit into `/etc/systemd/system/` with root-only perms.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1778")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
