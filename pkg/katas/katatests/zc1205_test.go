// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1205(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid ip neigh",
			input:    `ip neigh show`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid arp",
			input: `arp -a`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1205",
					Message: "Avoid `arp` — it is deprecated on modern Linux. Use `ip neigh` from iproute2 for neighbor table management.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1205")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
