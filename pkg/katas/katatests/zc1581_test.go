// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1581(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — ssh user@host",
			input:    `ssh user@host`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — ssh -o PubkeyAuthentication=yes",
			input:    `ssh -o PubkeyAuthentication=yes user@host`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — ssh -o PubkeyAuthentication=no",
			input: `ssh -o PubkeyAuthentication=no user@host`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1581",
					Message: "`ssh -o PubkeyAuthentication=no` forces password auth — weaker than key auth. Let the default preference pick.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — ssh -o PasswordAuthentication=yes",
			input: `ssh -o PasswordAuthentication=yes user@host`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1581",
					Message: "`ssh -o PasswordAuthentication=yes` forces password auth — weaker than key auth. Let the default preference pick.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — ssh -o PreferredAuthentications=password",
			input: `ssh -o PreferredAuthentications=password user@host`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1581",
					Message: "`ssh -o PreferredAuthentications=password` forces password auth — weaker than key auth. Let the default preference pick.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1581")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
