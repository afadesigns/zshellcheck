// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1559(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — ssh-copy-id user@host",
			input:    `ssh-copy-id user@host`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — ssh-copy-id -i ~/.ssh/id_ed25519.pub user@host",
			input:    `ssh-copy-id -i ~/.ssh/id_ed25519.pub user@host`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — ssh-copy-id -f user@host",
			input: `ssh-copy-id -f user@host`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1559",
					Message: "`ssh-copy-id -f` pushes a long-term credential without host-key verification. Verify the fingerprint out of band first.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — ssh-copy-id -o StrictHostKeyChecking=no user@host",
			input: `ssh-copy-id -o StrictHostKeyChecking=no user@host`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1559",
					Message: "`ssh-copy-id -o StrictHostKeyChecking=no` pushes a long-term credential without host-key verification. Verify the fingerprint out of band first.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1559")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
