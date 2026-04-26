// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1556(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — openssl enc -aes-256-gcm",
			input:    `openssl enc -aes-256-gcm -in file -out enc`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — openssl enc -chacha20-poly1305",
			input:    `openssl enc -chacha20-poly1305 -in file -out enc`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — openssl enc -des",
			input: `openssl enc -des -in file -out enc`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1556",
					Message: "`openssl enc -des` is a broken or deprecated cipher. Use `-aes-256-gcm` / `-chacha20-poly1305`, or `age` / `gpg` for files.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — openssl enc -rc4",
			input: `openssl enc -rc4 -in file -out enc`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1556",
					Message: "`openssl enc -rc4` is a broken or deprecated cipher. Use `-aes-256-gcm` / `-chacha20-poly1305`, or `age` / `gpg` for files.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — openssl enc -des-ede3-cbc",
			input: `openssl enc -des-ede3-cbc -in file -out enc`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1556",
					Message: "`openssl enc -des-ede3-cbc` is a broken or deprecated cipher. Use `-aes-256-gcm` / `-chacha20-poly1305`, or `age` / `gpg` for files.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1556")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
