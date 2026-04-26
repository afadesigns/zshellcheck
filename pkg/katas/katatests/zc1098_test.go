// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1098(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "no eval",
			input:    `echo hello`,
			expected: []katas.Violation{},
		},
		{
			name:     "eval without variables",
			input:    `eval "echo hello"`,
			expected: []katas.Violation{},
		},
		{
			name:  "eval with unquoted variable",
			input: `eval "ls $dir"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1098",
					Message: "Use the `(q)` flag (or `(qq)`, `(q-)`) when using variables in `eval` to prevent injection.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:     "eval with quoted variable",
			input:    `eval "ls ${(q)dir}"`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1098")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
