// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1063(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:  "fgrep usage",
			input: `fgrep 'literal' file`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1063",
					Message: "`fgrep` is deprecated. Use `grep -F` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:     "grep -F usage",
			input:    `grep -F 'literal' file`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1063")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
