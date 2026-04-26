// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1330(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid ZDOTDIR usage",
			input:    `echo $ZDOTDIR`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid INPUTRC usage",
			input: `echo $INPUTRC`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1330",
					Message: "Avoid `$INPUTRC` in Zsh — Zsh uses `bindkey` and ZLE, not readline. `INPUTRC` is Bash-specific.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1330")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
