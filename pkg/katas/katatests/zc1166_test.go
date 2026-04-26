// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1166(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid grep -i with file",
			input:    `grep -i pattern file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid grep -i in pipeline",
			input: `grep -i pattern`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1166",
					Message: "Use Zsh `(#i)` glob flag for case-insensitive matching instead of piping through `grep -i`. Example: `[[ $var == (#i)pattern ]]`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1166")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
