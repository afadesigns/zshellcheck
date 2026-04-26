// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1332(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid variable",
			input:    `echo $MYGLOB`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid GLOBIGNORE usage",
			input: `echo $GLOBIGNORE`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1332",
					Message: "Avoid `$GLOBIGNORE` in Zsh — use `setopt EXTENDED_GLOB` with `~` operator for glob exclusion.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1332")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
