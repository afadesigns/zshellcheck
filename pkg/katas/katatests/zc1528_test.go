// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1528(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — chage -M 90 alice",
			input:    `chage -M 90 alice`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — chage -l alice",
			input:    `chage -l alice`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — chage -M 99999 alice",
			input: `chage -M 99999 alice`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1528",
					Message: "`chage -M 99999` disables password aging — removes automatic lockout. Use a PAM profile instead of per-user chage.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — chage -E -1 alice",
			input: `chage -E -1 alice`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1528",
					Message: "`chage -E -1` disables password aging — removes automatic lockout. Use a PAM profile instead of per-user chage.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1528")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
