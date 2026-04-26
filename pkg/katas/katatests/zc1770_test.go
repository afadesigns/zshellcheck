// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1770(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `gpg --verify sig.asc` (default trust model)",
			input:    `gpg --verify sig.asc`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `gpg --trust-model pgp --verify sig.asc`",
			input:    `gpg --trust-model pgp --verify sig.asc`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `gpg --verify --always-trust sig.asc` (trailing form)",
			input: `gpg --verify --always-trust sig.asc`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1770",
					Message: "`gpg --always-trust` marks every imported key as fully trusted — a signature from an attacker-supplied key verifies cleanly. Drop the flag and pin the expected fingerprint, or assign trust via `gpg --edit-key KEYID trust`.",
					Line:    1,
					Column:  15,
				},
			},
		},
		{
			name:  "invalid — `gpg --verify --trust-model always sig.asc`",
			input: `gpg --verify --trust-model always sig.asc`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1770",
					Message: "`gpg --trust-model always` marks every imported key as fully trusted — a signature from an attacker-supplied key verifies cleanly. Drop the flag and pin the expected fingerprint, or assign trust via `gpg --edit-key KEYID trust`.",
					Line:    1,
					Column:  15,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1770")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
