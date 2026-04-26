// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1172(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid read -A",
			input:    `read -A arr`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid read -r",
			input:    `read -r line`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid read -a (Bash syntax)",
			input: `read -a arr`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1172",
					Message: "Use `read -A` instead of `read -a` in Zsh. The `-a` flag is Bash syntax; Zsh uses `-A` to read into arrays.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1172")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
