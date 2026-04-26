// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1250(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid gpg -b -d",
			input:    `gpg -b -d file.gpg`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid gpg without operation",
			input:    `gpg -k`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid gpg -d without -b",
			input: `gpg -d file.gpg`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1250",
					Message: "Use `gpg --batch` in scripts for non-interactive operation. Without `--batch`, gpg may prompt for passphrases or confirmations.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1250")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
