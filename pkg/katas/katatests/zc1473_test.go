// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1473(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — openssl req with -aes256",
			input:    `openssl req -newkey rsa:4096 -aes256 -keyout key.pem -out csr.pem`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — openssl x509 (not key-producing)",
			input:    `openssl x509 -in cert.pem -noout -subject`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — openssl req -nodes",
			input: `openssl req -newkey rsa:4096 -nodes -keyout key.pem -out csr.pem`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1473",
					Message: "`-nodes` writes the private key to disk unencrypted. Use `-aes256` (or an HSM/TPM) and keep the passphrase in a secrets store.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — openssl genrsa with -noenc",
			input: `openssl genrsa -noenc -out key.pem 4096`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1473",
					Message: "`-noenc` writes the private key to disk unencrypted. Use `-aes256` (or an HSM/TPM) and keep the passphrase in a secrets store.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1473")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
