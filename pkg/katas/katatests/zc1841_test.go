// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1841(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `curl --proxy-cacert /etc/ssl/proxy.pem https://api`",
			input:    `curl --proxy-cacert /etc/ssl/proxy.pem https://api`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `curl https://api` (no proxy flags)",
			input:    `curl https://api`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `curl --proxy-insecure https://api` (flag first)",
			input: `curl --proxy-insecure https://api`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1841",
					Message: "`curl --proxy-insecure` skips TLS verification on the proxy hop — an on-path attacker can present any cert and decrypt the tunnel (including `Authorization:` headers). Install the proxy CA and use `--proxy-cacert PATH`.",
					Line:    1,
					Column:  7,
				},
			},
		},
		{
			name:  "invalid — `curl https://api --proxy-insecure` (flag trailing)",
			input: `curl https://api --proxy-insecure`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1841",
					Message: "`curl --proxy-insecure` skips TLS verification on the proxy hop — an on-path attacker can present any cert and decrypt the tunnel (including `Authorization:` headers). Install the proxy CA and use `--proxy-cacert PATH`.",
					Line:    1,
					Column:  19,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1841")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
