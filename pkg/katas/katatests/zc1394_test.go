// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1394(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — echo $ZSH_NAME",
			input:    `echo $ZSH_NAME`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — echo $BASH",
			input: `echo $BASH`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1394",
					Message: "`$BASH` is Bash-only. Zsh exposes the interpreter name via `$ZSH_NAME` and the executable path indirectly via `$0`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1394")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
