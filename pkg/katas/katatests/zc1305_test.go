// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1305(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid words usage",
			input:    `echo $words`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid COMP_WORDS usage",
			input: `echo $COMP_WORDS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1305",
					Message: "Avoid `$COMP_WORDS` in Zsh — use `$words` array instead. `COMP_WORDS` is Bash completion-specific.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1305")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
