// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1186(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid unset -v",
			input:    `unset -v myvar`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid unset -f",
			input:    `unset -f myfunc`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid bare unset",
			input: `unset myvar`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1186",
					Message: "Use `unset -v name` for variables or `unset -f name` for functions. Bare `unset` is ambiguous about what is being removed.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1186")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
