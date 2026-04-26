// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1785(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `ufw default deny incoming`",
			input:    `ufw default deny incoming`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `ufw allow 22/tcp`",
			input:    `ufw allow 22/tcp`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `ufw default allow incoming`",
			input: `ufw default allow incoming`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1785",
					Message: "`ufw default allow incoming` flips the firewall baseline to accept every port that is not explicitly denied. Restore with `ufw default deny incoming` and add narrow `ufw allow <port>` rules for the services that must be reachable.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `ufw default allow` (direction omitted)",
			input: `ufw default allow`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1785",
					Message: "`ufw default allow incoming` flips the firewall baseline to accept every port that is not explicitly denied. Restore with `ufw default deny incoming` and add narrow `ufw allow <port>` rules for the services that must be reachable.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1785")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
