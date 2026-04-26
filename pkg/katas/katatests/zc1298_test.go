// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1298(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid funcstack usage",
			input:    `echo $funcstack`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid FUNCNAME usage",
			input: `echo $FUNCNAME`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1298",
					Message: "Avoid `$FUNCNAME` in Zsh — use `$funcstack` instead. `FUNCNAME` is Bash-specific.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1298")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
