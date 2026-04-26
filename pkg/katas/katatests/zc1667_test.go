// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1667(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — openssl enc with -pbkdf2",
			input:    `openssl enc -aes-256-cbc -pbkdf2 -iter 100000 -in file -out file.enc`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — openssl req (different subcommand)",
			input:    `openssl req -new -key key.pem -out csr.pem`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — openssl enc -aes-256-cbc without pbkdf2",
			input: `openssl enc -aes-256-cbc -salt -in file -out file.enc`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1667",
					Message: "`openssl enc` without `-pbkdf2` uses single-round EVP_BytesToKey (MD5) — add `-pbkdf2 -iter 100000`, or prefer `age` / `gpg --symmetric`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — openssl enc -aes-256-gcm (no pbkdf2, no AEAD support)",
			input: `openssl enc -aes-256-gcm -in file -out file.enc`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1667",
					Message: "`openssl enc` without `-pbkdf2` uses single-round EVP_BytesToKey (MD5) — add `-pbkdf2 -iter 100000`, or prefer `age` / `gpg --symmetric`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1667")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
