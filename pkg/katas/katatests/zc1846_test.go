// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1846(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `certbot renew` (default)",
			input:    `certbot renew`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `certbot certificates`",
			input:    `certbot certificates`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `certbot renew --force-renewal`",
			input: `certbot renew --force-renewal`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1846",
					Message: "`certbot renew --force-renewal` reissues regardless of expiry — in a cron it burns Let's Encrypt rate limits (50 certs per domain / 7 days). Drop the flag and let the 30-day gate work.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `certbot certonly --force-renewal -d example.com`",
			input: `certbot certonly --force-renewal -d example.com`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1846",
					Message: "`certbot certonly --force-renewal` reissues regardless of expiry — in a cron it burns Let's Encrypt rate limits (50 certs per domain / 7 days). Drop the flag and let the 30-day gate work.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1846")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
