// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1320(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid argv usage",
			input:    `echo $argv`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid BASH_ARGV usage",
			input: `echo $BASH_ARGV`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1320",
					Message: "Avoid `$BASH_ARGV` in Zsh — use `$argv` or `$@` for positional parameters. `BASH_ARGV` is Bash-specific.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1320")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
