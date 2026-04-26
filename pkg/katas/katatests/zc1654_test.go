// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1654(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — sysctl -p /etc/sysctl.conf",
			input:    `sysctl -p /etc/sysctl.conf`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — sysctl -p (default)",
			input:    `sysctl -p`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — sysctl -p /tmp/sysctl.conf",
			input: `sysctl -p /tmp/sysctl.conf`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1654",
					Message: "`sysctl -p /tmp/sysctl.conf` reads tunables from a world-traversable path — a concurrent local user can substitute the file. Keep configs under `/etc/sysctl.d/`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — sysctl -p /var/tmp/x",
			input: `sysctl -p /var/tmp/x.conf`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1654",
					Message: "`sysctl -p /var/tmp/x.conf` reads tunables from a world-traversable path — a concurrent local user can substitute the file. Keep configs under `/etc/sysctl.d/`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1654")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
