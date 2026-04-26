// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1064(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:  "type command",
			input: `type ls`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1064",
					Message: "Prefer `command -v` over `type`. `type` output is not stable/standard for checking command existence.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:     "command -v usage",
			input:    `command -v ls`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1064")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
