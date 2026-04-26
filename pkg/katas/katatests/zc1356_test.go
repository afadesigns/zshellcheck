// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1356(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — read -A",
			input:    `read -A arr`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — read -r line",
			input:    `read -r line`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — read -a (Bash syntax)",
			input: `read -a arr`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1356",
					Message: "Use `read -A` (uppercase) in Zsh to read into an array. `read -a` has different semantics in Zsh than in Bash.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1356")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
