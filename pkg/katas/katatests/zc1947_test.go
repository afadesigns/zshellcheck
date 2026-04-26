// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1947(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `ip xfrm state list`",
			input:    `ip xfrm state list`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `ip xfrm policy show`",
			input:    `ip xfrm policy show`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `ip xfrm state flush`",
			input: `ip xfrm state flush`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1947",
					Message: "`ip xfrm state flush` tears down every IPsec SA/policy — VPN tunnels drop, kernel stops encrypting, plaintext may leak during renegotiation. Scope via `ip xfrm state deleteall src $A dst $B`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `ip xfrm policy flush`",
			input: `ip xfrm policy flush`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1947",
					Message: "`ip xfrm policy flush` tears down every IPsec SA/policy — VPN tunnels drop, kernel stops encrypting, plaintext may leak during renegotiation. Scope via `ip xfrm policy deleteall src $A dst $B`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1947")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
