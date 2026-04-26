// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1487(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — history (list)",
			input:    `history`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — history 10",
			input:    `history 10`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — history -c",
			input: `history -c`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1487",
					Message: "`history -c` is a Bash-ism for clearing history — does nothing in Zsh and is a classic post-compromise tactic elsewhere.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — history -d 1",
			input: `history -d 1`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1487",
					Message: "`history -d` is a Bash-ism for clearing history — does nothing in Zsh and is a classic post-compromise tactic elsewhere.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1487")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
