// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1344(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — find without -size",
			input:    `find . -type f`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — find -size +10M",
			input: `find . -size +10M`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1344",
					Message: "Use Zsh `*(L±N[kmp])` glob qualifier instead of `find -size`. Unit suffixes: `k` kilobytes, `m` megabytes, `p` 512-byte blocks.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — find -size -1k",
			input: `find . -size -1k`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1344",
					Message: "Use Zsh `*(L±N[kmp])` glob qualifier instead of `find -size`. Unit suffixes: `k` kilobytes, `m` megabytes, `p` 512-byte blocks.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1344")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
