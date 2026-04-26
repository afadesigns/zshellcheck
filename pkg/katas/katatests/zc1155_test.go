// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1155(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid whence -a",
			input:    `whence -a git`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid which -a",
			input: `which -a git`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1155",
					Message: "Use `whence -a` instead of `which -a`. Zsh `whence` is a reliable builtin for listing all command locations.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1155")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
