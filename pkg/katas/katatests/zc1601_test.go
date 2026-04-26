// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1601(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — ethtool wol d (disable)",
			input:    `ethtool -s eth0 wol d`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — ethtool setting different knob",
			input:    `ethtool -s eth0 autoneg on`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — ethtool -s eth0 wol g",
			input: `ethtool -s eth0 wol g`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1601",
					Message: "`ethtool -s eth0 wol g` enables Wake-on-LAN — the NIC powers the host on before firewall rules load. Keep `wol d` unless a documented operational need requires g.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — ethtool -s $IF wol ubg",
			input: `ethtool -s $IF wol ubg`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1601",
					Message: "`ethtool -s $IF wol ubg` enables Wake-on-LAN — the NIC powers the host on before firewall rules load. Keep `wol d` unless a documented operational need requires ubg.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1601")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
