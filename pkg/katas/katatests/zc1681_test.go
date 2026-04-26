// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1681(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — tar -xf archive.tar",
			input:    `tar -xf archive.tar`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — tar -xvzf archive.tgz (no P)",
			input:    `tar -xvzf archive.tgz`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — tar -xPf archive.tar (short-flag cluster)",
			input: `tar -xPf archive.tar`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1681",
					Message: "`tar -xPf` keeps absolute paths during extraction — an untrusted archive can overwrite `/etc/cron.d`, `/usr/local/bin`, etc. Drop the flag and extract with `-C <scratch-dir>` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — tar xf archive.tar -P",
			input: `tar xf archive.tar -P`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1681",
					Message: "`tar -P` keeps absolute paths during extraction — an untrusted archive can overwrite `/etc/cron.d`, `/usr/local/bin`, etc. Drop the flag and extract with `-C <scratch-dir>` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1681")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
