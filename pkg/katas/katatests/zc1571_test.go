// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1571(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — chronyc makestep",
			input:    `chronyc makestep`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — ntpdate pool.ntp.org",
			input: `ntpdate pool.ntp.org`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1571",
					Message: "`ntpdate` is deprecated and races any running chrony/timesyncd. Use `chronyc makestep` or `systemctl restart systemd-timesyncd`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — sntp -sS pool.ntp.org",
			input: `sntp -sS pool.ntp.org`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1571",
					Message: "`sntp` is deprecated and races any running chrony/timesyncd. Use `chronyc makestep` or `systemctl restart systemd-timesyncd`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1571")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
