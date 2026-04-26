// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1334(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid whence -p usage",
			input:    `whence -p git`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid type without -p",
			input:    `type git`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid type -p usage",
			input: `type -p git`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1334",
					Message: "Avoid `type -p` in Zsh — use `whence -p` to get the command path instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1334")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
