// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1439(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — sysctl vm.swappiness=10",
			input:    `sysctl vm.swappiness=10`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — sysctl ip_forward=1",
			input: `sysctl net.ipv4.ip_forward=1`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1439",
					Message: "Enabling `ip_forward` turns the host into a router. Verify firewall posture (iptables/nftables) and persist the setting in `/etc/sysctl.d/`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1439")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
