// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1578(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — ssh-keygen -t ed25519 -f key",
			input:    `ssh-keygen -t ed25519 -f key`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — ssh-keygen -t rsa -b 4096",
			input:    `ssh-keygen -t rsa -b 4096`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — ssh-keygen -t rsa -b 1024",
			input: `ssh-keygen -t rsa -b 1024`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1578",
					Message: "`ssh-keygen -b 1024` — RSA below 2048 bits is rejected by modern OpenSSH. Use `-t ed25519` or `-b 4096`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — ssh-keygen -t dsa",
			input: `ssh-keygen -t dsa -f key`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1578",
					Message: "`ssh-keygen -t dsa` — DSA removed from OpenSSH 9.8. Use `-t ed25519`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1578")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
