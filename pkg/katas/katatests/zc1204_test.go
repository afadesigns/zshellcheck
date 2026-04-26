// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1204(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid ip route",
			input:    `ip route show`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid route",
			input: `route -n`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1204",
					Message: "Avoid `route` — it is deprecated on modern Linux. Use `ip route` from iproute2 for consistent routing management.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1204")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
