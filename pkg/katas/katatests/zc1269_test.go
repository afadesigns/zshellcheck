// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1269(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid pgrep usage",
			input:    `pgrep -f myprocess`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid ps with no grep-related args",
			input:    `ps -p 1234`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid ps aux for process search",
			input: `ps aux`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1269",
					Message: "Use `pgrep` instead of `ps aux | grep`. `pgrep` is purpose-built for process searching and doesn't match itself.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid ps -ef for process search",
			input: `ps -ef`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1269",
					Message: "Use `pgrep` instead of `ps -ef | grep`. `pgrep` is purpose-built for process searching and doesn't match itself.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid ps -e for process search",
			input: `ps -e`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1269",
					Message: "Use `pgrep` instead of `ps -e | grep`. `pgrep` is purpose-built for process searching and doesn't match itself.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1269")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
