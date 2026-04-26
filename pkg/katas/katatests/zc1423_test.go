// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1423(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — iptables -A",
			input:    `iptables -A INPUT -p tcp --dport 22 -j ACCEPT`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — iptables -F",
			input: `iptables -F`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1423",
					Message: "Flushing firewall rules with `-F` removes every rule — risk of locking yourself out of remote hosts. Save + use rollback mechanism.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — nft flush ruleset",
			input: `nft flush ruleset`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1423",
					Message: "`nft flush ruleset` clears every firewall table — risk of locking yourself out of remote hosts. Save + use rollback mechanism.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1423")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
