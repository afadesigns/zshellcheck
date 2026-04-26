// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1140(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid hash -r",
			input:    `hash -r`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid hash for existence check",
			input: `hash git`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1140",
					Message: "Use `command -v cmd` instead of `hash cmd` for command existence checks. `command -v` provides clearer semantics in Zsh.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1140")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
