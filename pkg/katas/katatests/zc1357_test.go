// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1357(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — printf %s",
			input:    `printf '%s\n' "$line"`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — printf %q",
			input: `printf '%q' "$v"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1357",
					Message: "Use Zsh `${(q)var}` for shell-quoting instead of `printf '%q'`. Variants: `${(qq)}`, `${(qqq)}`, `${(qqqq)}` for different quote styles.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1357")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
