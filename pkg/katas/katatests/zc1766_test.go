// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1766(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `memcached -l 127.0.0.1`",
			input:    `memcached -l 127.0.0.1`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `memcached -l 10.0.0.5 -p 11211`",
			input:    `memcached -l 10.0.0.5 -p 11211`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `memcached -l 0.0.0.0`",
			input: `memcached -l 0.0.0.0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1766",
					Message: "`memcached -l 0.0.0.0` exposes the unauthenticated cache to every interface on the host. Bind to `127.0.0.1` or a private-network IP and firewall the port.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `memcached -l0.0.0.0` (joined form)",
			input: `memcached -l0.0.0.0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1766",
					Message: "`memcached -l0.0.0.0` exposes the unauthenticated cache to every interface on the host. Bind to `127.0.0.1` or a private-network IP and firewall the port.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1766")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
