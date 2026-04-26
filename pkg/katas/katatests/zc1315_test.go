// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1315(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid emulate usage",
			input:    `emulate -L sh`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid BASH_COMPAT usage",
			input: `echo $BASH_COMPAT`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1315",
					Message: "Avoid `$BASH_COMPAT` in Zsh — use `emulate` for shell compatibility mode instead.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1315")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
