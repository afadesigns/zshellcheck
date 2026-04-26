// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1314(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid non-Bash variable",
			input:    `echo $MY_PATH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid BASH_LOADABLES_PATH",
			input: `echo $BASH_LOADABLES_PATH`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1314",
					Message: "Avoid `$BASH_LOADABLES_PATH` in Zsh — it is undefined. Use `zmodload` with full module names.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1314")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
