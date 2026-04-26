// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1558(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — usermod -aG audio alice",
			input:    `usermod -aG audio alice`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — usermod -aG wheel alice",
			input: `usermod -aG wheel alice`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1558",
					Message: "Adding user to `wheel` grants persistent admin-level access — use a scoped sudoers.d drop-in via configuration management.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — usermod -aG docker alice",
			input: `usermod -aG docker alice`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1558",
					Message: "Adding user to `docker` grants persistent admin-level access — use a scoped sudoers.d drop-in via configuration management.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — gpasswd -a alice sudo",
			input: `gpasswd -a alice sudo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1558",
					Message: "Adding user to `sudo` grants persistent admin-level access — use a scoped sudoers.d drop-in via configuration management.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — usermod -aG audio,wheel alice (mixed)",
			input: `usermod -aG audio,wheel alice`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1558",
					Message: "Adding user to `wheel` grants persistent admin-level access — use a scoped sudoers.d drop-in via configuration management.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1558")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
