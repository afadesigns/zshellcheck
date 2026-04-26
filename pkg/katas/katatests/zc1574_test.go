// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1574(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — git config credential.helper libsecret",
			input:    `git config credential.helper libsecret`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — git config credential.helper 'cache --timeout=3600'",
			input:    `git config credential.helper 'cache --timeout=3600'`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — git config credential.helper store",
			input: `git config credential.helper store`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1574",
					Message: "`git credential.helper store` saves credentials in plaintext — backups leak the token. Use platform helper (manager-core / libsecret) or `cache --timeout=<sec>`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — git config --global credential.helper store",
			input: `git config --global credential.helper store`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1574",
					Message: "`git credential.helper store` saves credentials in plaintext — backups leak the token. Use platform helper (manager-core / libsecret) or `cache --timeout=<sec>`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1574")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
