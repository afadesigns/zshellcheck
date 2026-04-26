// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1125(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid grep -q with file",
			input:    `grep -q pattern file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid grep without -q",
			input:    `grep pattern`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid grep -q in pipeline",
			input: `grep -q pattern`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1125",
					Message: "Use `[[ $var =~ pattern ]]` or `[[ $var == *pattern* ]]` instead of piping through `grep -q`. Zsh pattern matching avoids spawning external processes.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1125")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
