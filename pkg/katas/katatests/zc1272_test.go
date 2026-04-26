// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1272(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid install usage",
			input:    `install -m 0755 mybin /usr/local/bin`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid cp to non-system dir",
			input:    `cp file.txt /home/user/`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid cp to /usr/local/bin",
			input: `cp mybin /usr/local/bin`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1272",
					Message: "Use `install -m 0755` instead of `cp` to system directories. `install` sets permissions atomically.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid cp to /usr/bin",
			input: `cp mybin /usr/bin`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1272",
					Message: "Use `install -m 0755` instead of `cp` to system directories. `install` sets permissions atomically.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1272")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
