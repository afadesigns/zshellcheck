// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1974(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `ipset list`",
			input:    `ipset list`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `ipset destroy blocklist` (targeted name)",
			input:    `ipset destroy blocklist`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `ipset flush`",
			input: `ipset flush`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1974",
					Message: "`ipset flush` drops named IP sets wholesale — iptables/nft rules that reference them fall through to the default policy (block-list empty, allow-list gone). Target by name; reload atomically via `ipset restore -! < snapshot`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `ipset destroy` (no arg, wipes every set)",
			input: `ipset destroy`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1974",
					Message: "`ipset destroy` drops named IP sets wholesale — iptables/nft rules that reference them fall through to the default policy (block-list empty, allow-list gone). Target by name; reload atomically via `ipset restore -! < snapshot`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1974")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
