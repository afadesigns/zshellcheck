// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1760(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `openssl rand -hex 32`",
			input:    `openssl rand -hex 32`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `openssl rand -hex 16` (borderline but accepted)",
			input:    `openssl rand -hex 16`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `openssl rand 24` (no encoding flag)",
			input:    `openssl rand 24`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `openssl rand -hex 8`",
			input: `openssl rand -hex 8`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1760",
					Message: "`openssl rand -hex 8` produces a sub-128-bit value — brute-forceable offline. Use `-hex 32` for secrets / long-lived tokens.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `openssl rand -base64 12`",
			input: `openssl rand -base64 12`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1760",
					Message: "`openssl rand -base64 12` produces a sub-128-bit value — brute-forceable offline. Use `-hex 32` for secrets / long-lived tokens.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1760")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
