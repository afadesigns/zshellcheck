// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1309(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid non-Bash variable",
			input:    `echo $MY_COMMAND`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid BASH_COMMAND usage",
			input: `echo $BASH_COMMAND`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1309",
					Message: "Avoid `$BASH_COMMAND` in Zsh — it is undefined. Use `$ZSH_DEBUG_CMD` in debug traps if needed.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1309")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
