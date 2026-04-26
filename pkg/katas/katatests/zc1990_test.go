// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1990(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `openssl passwd -5 $PASS` (SHA-256-crypt)",
			input:    `openssl passwd -5 $PASS`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `openssl passwd -6 $PASS` (SHA-512-crypt)",
			input:    `openssl passwd -6 $PASS`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `openssl passwd -crypt $PASS`",
			input: `openssl passwd -crypt $PASS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1990",
					Message: "`openssl passwd -crypt` emits a broken hash format — DES/MD5 variants crack on a laptop. Use `-5` / `-6` or a KDF-based hasher (`mkpasswd -m yescrypt`, `htpasswd -B`, `argon2`).",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `openssl passwd -1 $PASS`",
			input: `openssl passwd -1 $PASS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1990",
					Message: "`openssl passwd -1` emits a broken hash format — DES/MD5 variants crack on a laptop. Use `-5` / `-6` or a KDF-based hasher (`mkpasswd -m yescrypt`, `htpasswd -B`, `argon2`).",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `openssl passwd -apr1 $PASS`",
			input: `openssl passwd -apr1 $PASS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1990",
					Message: "`openssl passwd -apr1` emits a broken hash format — DES/MD5 variants crack on a laptop. Use `-5` / `-6` or a KDF-based hasher (`mkpasswd -m yescrypt`, `htpasswd -B`, `argon2`).",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1990")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
