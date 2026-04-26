// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1698(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — fail2ban-client status",
			input:    `fail2ban-client status`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — fail2ban-client set sshd unbanip scoped",
			input:    `fail2ban-client set sshd unbanip 1.2.3.4`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — fail2ban-client unban --all",
			input: `fail2ban-client unban --all`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1698",
					Message: "`fail2ban-client unban --all` wipes every active brute-force ban — attacker IPs regain access. Target individual IPs with `set <jail> unbanip <ip>` or reload a single jail.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — fail2ban-client stop",
			input: `fail2ban-client stop`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1698",
					Message: "`fail2ban-client stop` wipes every active brute-force ban — attacker IPs regain access. Target individual IPs with `set <jail> unbanip <ip>` or reload a single jail.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1698")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
