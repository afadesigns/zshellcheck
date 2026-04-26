// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1165(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid awk with complex script",
			input:    `awk '{sum+=$1} END{print sum}'`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid awk with file",
			input:    `awk '{print $1}' file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid simple awk print $1",
			input: `awk '{print $1}'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1165",
					Message: "Use Zsh parameter expansion (`${var%% *}` or `${var##* }`) instead of `awk '{print $1}'` for simple field extraction without spawning awk.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1165")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
