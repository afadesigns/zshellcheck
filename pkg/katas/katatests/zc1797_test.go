// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1797(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `ip link set eth0 up`",
			input:    `ip link set eth0 up`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `ip addr show` (read only)",
			input:    `ip addr show`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `ip link set eth0 down`",
			input: `ip link set eth0 down`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1797",
					Message: "`ip link set … down` disables a network interface — if it carries the SSH session, the script cuts itself off. Schedule a rollback via `systemd-run --on-active=30s ip link set … up` or stage via `nmcli`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `ifdown eth0`",
			input: `ifdown eth0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1797",
					Message: "`ifdown eth0` disables a network interface — if it carries the SSH session, the script cuts itself off. Schedule a rollback via `systemd-run --on-active=30s ip link set … up` or stage via `nmcli`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1797")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
