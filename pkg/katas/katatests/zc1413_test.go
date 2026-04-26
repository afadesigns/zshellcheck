// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1413(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — hash -r reset",
			input:    `hash -r`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — hash -t",
			input: `hash -t ls`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1413",
					Message: "Use `whence -p cmd` (Zsh) instead of `hash -t cmd`. `whence -p` always returns the absolute path, regardless of hash state.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1413")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
