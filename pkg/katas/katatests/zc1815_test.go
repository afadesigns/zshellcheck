// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1815(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `systemctl status NetworkManager` (read only)",
			input:    `systemctl status NetworkManager`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `systemctl restart nginx`",
			input:    `systemctl restart nginx`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `systemctl restart NetworkManager`",
			input: `systemctl restart NetworkManager`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1815",
					Message: "`systemctl restart NetworkManager` drops every connection the manager supervises — the SSH session freezes. Use `nmcli connection reload` / `networkctl reload`, or a `systemd-run --on-active=30s` rollback.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `systemctl restart systemd-networkd.service`",
			input: `systemctl restart systemd-networkd.service`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1815",
					Message: "`systemctl restart systemd-networkd.service` drops every connection the manager supervises — the SSH session freezes. Use `nmcli connection reload` / `networkctl reload`, or a `systemd-run --on-active=30s` rollback.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1815")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
