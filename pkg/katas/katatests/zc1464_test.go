// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1464(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — iptables -A INPUT -j DROP",
			input:    `iptables -A INPUT -j DROP`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — iptables -P INPUT DROP",
			input:    `iptables -P INPUT DROP`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — iptables-save > backup",
			input:    `iptables -S`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — iptables -F (flush all)",
			input: `iptables -F`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1464",
					Message: "Firewall hardening weakened (flushing all firewall rules). Keep default-drop and use atomic reload.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — iptables -P INPUT ACCEPT",
			input: `iptables -P INPUT ACCEPT`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1464",
					Message: "Firewall hardening weakened (default-ACCEPT policy on INPUT chain). Keep default-drop and use atomic reload.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — ip6tables -F",
			input: `ip6tables -F`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1464",
					Message: "Firewall hardening weakened (flushing all firewall rules). Keep default-drop and use atomic reload.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1464")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
