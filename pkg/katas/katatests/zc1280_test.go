// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1280(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid Zsh :e modifier",
			input:    `echo ${file:e}`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid cut with different delimiter",
			input:    `cut -d: -f1`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid cut to extract extension",
			input: `cut -d. -f2`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1280",
					Message: "Use Zsh parameter expansion `${var:e}` to extract the file extension instead of `cut -d. -f2`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1280")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
