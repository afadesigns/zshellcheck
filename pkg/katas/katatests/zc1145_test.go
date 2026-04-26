// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1145(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid tr with complex set",
			input:    `tr -d '[:space:]'`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid tr -d simple char",
			input: `tr -d ' '`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1145",
					Message: "Use `${var//char/}` instead of piping through `tr -d`. Parameter expansion is faster for simple character deletion.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1145")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
