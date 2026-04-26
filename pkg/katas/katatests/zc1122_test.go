// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1122(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:  "invalid whoami",
			input: `whoami`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1122",
					Message: "Use `$USER` instead of `whoami`. Zsh maintains `$USER` as a built-in variable, avoiding an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:     "valid other identifier",
			input:    `mycommand`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1122")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
