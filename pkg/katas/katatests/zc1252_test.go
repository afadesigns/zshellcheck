// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1252(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid getent",
			input:    `getent passwd root`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid cat /etc/passwd",
			input: `cat /etc/passwd`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1252",
					Message: "Use `getent` instead of `cat /etc/passwd`. `getent` queries all NSS sources including LDAP and SSSD.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1252")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
