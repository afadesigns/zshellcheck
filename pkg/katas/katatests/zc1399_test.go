// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1399(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — kill signal pid",
			input:    `kill -TERM 1234`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — kill -l",
			input: `kill -l`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1399",
					Message: "Use Zsh `print -l $signals` (after `zmodload zsh/parameter`) instead of `kill -l` for listing signal names.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1399")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
