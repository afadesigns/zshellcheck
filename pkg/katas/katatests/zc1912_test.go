// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1912(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `dhclient eth0` (renew)",
			input:    `dhclient eth0`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `dhcpcd --rebind eth0`",
			input:    `dhcpcd --rebind eth0`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `dhclient -r eth0`",
			input: `dhclient -r eth0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1912",
					Message: "`dhclient -r` drops the DHCP lease — SSH session cuts, VPC reachability stalls. Pair with a re-acquire (`dhclient -1`/`nmcli device reapply`), or schedule via `systemd-run --on-active=`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `dhcpcd -k eth0`",
			input: `dhcpcd -k eth0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1912",
					Message: "`dhcpcd -k` drops the DHCP lease — SSH session cuts, VPC reachability stalls. Pair with a re-acquire (`dhclient -1`/`nmcli device reapply`), or schedule via `systemd-run --on-active=`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1912")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
