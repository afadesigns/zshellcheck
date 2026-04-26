// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1848(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `ssh -o CheckHostIP=yes host`",
			input:    `ssh -o CheckHostIP=yes host`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `ssh host` (default)",
			input:    `ssh host`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `ssh -o CheckHostIP=no host` (split form)",
			input: `ssh -o CheckHostIP=no host`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1848",
					Message: "`ssh -o CheckHostIP=no` silences the IP-mismatch warning for known hosts — a DNS-spoof + leaked host-key attack goes undetected. Leave the default, or use `HostKeyAlias` for load-balanced pools.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `ssh -oCheckHostIP=no host` (attached form)",
			input: `ssh -oCheckHostIP=no host`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1848",
					Message: "`ssh -o CheckHostIP=no` silences the IP-mismatch warning for known hosts — a DNS-spoof + leaked host-key attack goes undetected. Leave the default, or use `HostKeyAlias` for load-balanced pools.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1848")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
