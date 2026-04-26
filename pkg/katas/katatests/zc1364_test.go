// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1364(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — cut -f (field, different kata)",
			input:    `cut -f 2 file`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — cut -c",
			input: `cut -c 1-5 file`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1364",
					Message: "Use Zsh `${var:pos:len}` for character ranges instead of `cut -c`. Parameter expansion is in-shell and zero-indexed.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — cut -c attached",
			input: `cut -c1-5 file`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1364",
					Message: "Use Zsh `${var:pos:len}` for character ranges instead of `cut -c`. Parameter expansion is in-shell and zero-indexed.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1364")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
