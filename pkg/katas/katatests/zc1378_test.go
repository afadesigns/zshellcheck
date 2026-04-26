// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1378(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — echo $dirstack (lowercase)",
			input:    `echo $dirstack`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — echo $DIRSTACK",
			input: `echo $DIRSTACK`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1378",
					Message: "Use lowercase `$dirstack` in Zsh — uppercase `$DIRSTACK` is Bash-only.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1378")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
