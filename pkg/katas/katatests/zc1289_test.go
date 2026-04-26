// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1289(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid sort with numeric and unique",
			input:    `sort -n -u file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid sort without unique",
			input:    `sort file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid sort -u alone",
			input: `sort -u file.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1289",
					Message: "Use Zsh `${(u)array}` for unique elements instead of `sort -u`. The `(u)` flag preserves order.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1289")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
