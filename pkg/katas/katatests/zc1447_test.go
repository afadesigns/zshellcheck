// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1447(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — ip addr",
			input:    `ip addr show`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — ifconfig",
			input: `ifconfig eth0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1447",
					Message: "`ifconfig` is deprecated. Use `ip addr` / `ip link` / `ip route` from iproute2.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — netstat",
			input: `netstat -tuln`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1447",
					Message: "`netstat` is deprecated. Use `ss` from iproute2 (same flags, faster output).",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1447")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
