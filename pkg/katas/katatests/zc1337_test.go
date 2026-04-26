// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1337(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid print usage",
			input:    `print -l hello`,
			expected: []katas.Violation{},
		},
		{
			name:  "fold usage",
			input: `fold -w 80 file.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1337",
					Message: "Consider Zsh `$COLUMNS` and `print` for text wrapping instead of `fold`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1337")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
