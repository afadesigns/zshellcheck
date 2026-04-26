// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1564(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — date (read)",
			input:    `date`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — timedatectl status",
			input:    `timedatectl status`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — date -s 2025-01-01",
			input: `date -s 2025-01-01`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1564",
					Message: "`date -s` sets the wall clock manually — breaks TLS certs, cron catch-up, and systemd timer math. Use timesyncd/chrony/ntpd.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — timedatectl set-time",
			input: `timedatectl set-time 2025-01-01`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1564",
					Message: "`timedatectl set-time` sets the wall clock manually — breaks TLS certs, cron catch-up, and systemd timer math. Use timesyncd/chrony/ntpd.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — hwclock -w",
			input: `hwclock -w`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1564",
					Message: "`hwclock -w` sets the wall clock manually — breaks TLS certs, cron catch-up, and systemd timer math. Use timesyncd/chrony/ntpd.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1564")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
