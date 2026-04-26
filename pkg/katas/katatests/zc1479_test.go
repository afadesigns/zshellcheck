// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1479(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — ssh user@host",
			input:    `ssh user@host`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — ssh -o StrictHostKeyChecking=accept-new",
			input:    `ssh -o StrictHostKeyChecking=accept-new user@host`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — ssh -o StrictHostKeyChecking=no",
			input: `ssh -o StrictHostKeyChecking=no user@host`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1479",
					Message: "`StrictHostKeyChecking=no` disables SSH host-key verification — first MITM owns the connection. Pin the fingerprint in known_hosts instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — scp -oStrictHostKeyChecking=no (joined)",
			input: `scp -oStrictHostKeyChecking=no file user@host:`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1479",
					Message: "`StrictHostKeyChecking=no` disables SSH host-key verification — first MITM owns the connection. Pin the fingerprint in known_hosts instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — ssh -oUserKnownHostsFile=/dev/null",
			input: `ssh -oUserKnownHostsFile=/dev/null user@host`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1479",
					Message: "`UserKnownHostsFile=/dev/null` disables SSH host-key verification — first MITM owns the connection. Pin the fingerprint in known_hosts instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1479")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
