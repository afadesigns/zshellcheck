// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1121(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid hostname -f",
			input:    `hostname -f`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid hostname -s",
			input:    `hostname -s`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid simple hostname",
			input: `hostname myhost`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1121",
					Message: "Use `$HOST` instead of `hostname`. Zsh maintains `$HOST` as a built-in variable, avoiding an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1121")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
