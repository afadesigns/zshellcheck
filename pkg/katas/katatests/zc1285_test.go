// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1285(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid sort with flags",
			input:    `sort -n file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid sort with key",
			input:    `sort -k 2`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid sort with reverse flag",
			input:    `sort -r file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:  "sort with single file argument",
			input: `sort data.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1285",
					Message: "Use Zsh `${(o)array}` for sorting instead of piping to `sort`. The `(o)` flag sorts in-shell.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1285")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
