// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1142(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid single grep",
			input:    `grep -E "foo|bar" file`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid grep | grep",
			input: `grep foo file | grep bar`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1142",
					Message: "Avoid chaining `grep | grep`. Combine into a single `grep -E` with alternation or use `awk` for multi-pattern matching to reduce pipeline processes.",
					Line:    1,
					Column:  15,
				},
			},
		},
		{
			name:     "non-pipe operator",
			input:    `echo hello && echo world`,
			expected: []katas.Violation{},
		},
		{
			name:     "pipe but not grep",
			input:    `cat file | sort`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1142")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
