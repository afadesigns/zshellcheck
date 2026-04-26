// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1723(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `gpg --list-keys`",
			input:    `gpg --list-keys`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `gpg KEYID --export-secret-keys`",
			input:    `gpg KEYID --export-secret-keys`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `gpg --delete-secret-keys KEYID` (leading-flag form)",
			input: `gpg --delete-secret-keys KEYID`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1723",
					Message: "`gpg --delete-secret-keys` permanently destroys keyring entries — no recovery without a separate backup. Export with `gpg --export-secret-keys --armor KEYID` first; never pair this flag with `--batch --yes`.",
					Line:    1,
					Column:  6,
				},
			},
		},
		{
			name:  "invalid — `gpg --batch --yes --delete-secret-and-public-keys KEYID`",
			input: `gpg --batch --yes --delete-secret-and-public-keys KEYID`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1723",
					Message: "`gpg --delete-secret-and-public-keys` permanently destroys keyring entries — no recovery without a separate backup. Export with `gpg --export-secret-keys --armor KEYID` first; never pair this flag with `--batch --yes`.",
					Line:    1,
					Column:  20,
				},
			},
		},
		{
			name:  "invalid — `gpg KEYID --delete-key` (trailing-flag form)",
			input: `gpg KEYID --delete-key`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1723",
					Message: "`gpg --delete-key` permanently destroys keyring entries — no recovery without a separate backup. Export with `gpg --export-secret-keys --armor KEYID` first; never pair this flag with `--batch --yes`.",
					Line:    1,
					Column:  12,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1723")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
