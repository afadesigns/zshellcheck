// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1291(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid sort -r with -n",
			input:    `sort -r -n file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid sort without -r",
			input:    `sort file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid sort -r alone",
			input: `sort -r file.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1291",
					Message: "Use Zsh `${(O)array}` for reverse sorting instead of `sort -r`. The `(O)` flag sorts descending in-shell.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1291")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
