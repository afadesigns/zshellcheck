// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1279(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid realpath usage",
			input:    `realpath /some/path`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid readlink without -f",
			input:    `readlink /some/symlink`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid readlink -f",
			input: `readlink -f /some/path`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1279",
					Message: "Use `realpath` instead of `readlink -f`. `realpath` is more portable, especially on macOS.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1279")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
