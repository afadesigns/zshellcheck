// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1317(t *testing.T) {
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
			name:  "invalid BASH_ENV usage",
			input: `echo $BASH_ENV`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1317",
					Message: "Avoid `$BASH_ENV` in Zsh — use `$ZDOTDIR` for Zsh startup file locations instead.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1317")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
