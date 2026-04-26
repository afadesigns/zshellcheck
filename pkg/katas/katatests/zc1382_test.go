// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1382(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — echo $BUFFER (Zsh ZLE)",
			input:    `echo $BUFFER`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — echo $READLINE_LINE",
			input: `echo $READLINE_LINE`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1382",
					Message: "Bash `$READLINE_*` vars do not exist in Zsh. Inside ZLE widgets use `$BUFFER`, `$CURSOR`, `$MARK`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — echo $READLINE_POINT",
			input: `echo $READLINE_POINT`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1382",
					Message: "Bash `$READLINE_*` vars do not exist in Zsh. Inside ZLE widgets use `$BUFFER`, `$CURSOR`, `$MARK`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1382")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
