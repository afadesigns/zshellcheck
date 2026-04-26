// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1485(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — openssl s_client -tls1_3",
			input:    `openssl s_client -tls1_3 -connect host:443`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — openssl x509 (not s_client/s_server)",
			input:    `openssl x509 -ssl3 -noout`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — openssl s_client -ssl3",
			input: `openssl s_client -ssl3 -connect host:443`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1485",
					Message: "`openssl s_client -ssl3` forces a legacy / disabled TLS version (downgrade-attack surface). Update the remote instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — openssl s_client -tls1",
			input: `openssl s_client -tls1 -connect host:443`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1485",
					Message: "`openssl s_client -tls1` forces a legacy / disabled TLS version (downgrade-attack surface). Update the remote instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — openssl s_server -tls1_1",
			input: `openssl s_server -tls1_1 -cert cert.pem`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1485",
					Message: "`openssl s_server -tls1_1` forces a legacy / disabled TLS version (downgrade-attack surface). Update the remote instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1485")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
