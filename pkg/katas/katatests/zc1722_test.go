// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1722(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `ssh-keyscan HOST` (no redirect; fingerprint check separately)",
			input:    `ssh-keyscan HOST`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `ssh-keyscan HOST > /tmp/scan.tmp` (not known_hosts)",
			input:    `ssh-keyscan HOST > /tmp/scan.tmp`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `ssh-keyscan HOST >> ~/.ssh/known_hosts`",
			input: `ssh-keyscan HOST >> ~/.ssh/known_hosts`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1722",
					Message: "`ssh-keyscan ... >> ~/.ssh/known_hosts` accepts the first-served host key without verifying its fingerprint. Pipe to `ssh-keygen -lf -` and assert the fingerprint first.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `ssh-keyscan -H HOST > known_hosts`",
			input: `ssh-keyscan -H HOST > known_hosts`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1722",
					Message: "`ssh-keyscan ... > known_hosts` accepts the first-served host key without verifying its fingerprint. Pipe to `ssh-keygen -lf -` and assert the fingerprint first.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1722")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
