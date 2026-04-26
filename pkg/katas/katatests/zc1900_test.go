// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1900(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `curl -L https://api/resource`",
			input:    `curl -L https://api/resource`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `curl -u user:pass https://api/resource` (no location)",
			input:    `curl -u user:pass https://api/resource`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `curl --location-trusted https://api` (leading)",
			input: `curl --location-trusted https://api`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1900",
					Message: "`curl --location-trusted` replays `Authorization`, cookies, and `-u user:pass` on every redirect — a 302 to attacker-controlled host leaks the token. Drop the flag; verify final hostname before sending secrets.",
					Line:    1,
					Column:  7,
				},
			},
		},
		{
			name:  "invalid — `curl -u user:pass --location-trusted https://api` (trailing)",
			input: `curl -u user:pass --location-trusted https://api`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1900",
					Message: "`curl --location-trusted` replays `Authorization`, cookies, and `-u user:pass` on every redirect — a 302 to attacker-controlled host leaks the token. Drop the flag; verify final hostname before sending secrets.",
					Line:    1,
					Column:  20,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1900")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
