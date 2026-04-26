// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1243(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid grep -lZ",
			input:    `grep -lZ pattern dir`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid grep without -l",
			input:    `grep pattern file`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid grep -l without -Z",
			input: `grep -l pattern dir`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1243",
					Message: "Use `grep -lZ` instead of `grep -l` for null-terminated file lists. Pair with `xargs -0` to safely handle filenames with special characters.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1243")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
