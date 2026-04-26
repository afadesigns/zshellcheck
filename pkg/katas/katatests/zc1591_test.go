// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1591(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — printf with specific format",
			input:    `printf '%-20s %d\n' "${pairs[@]}"`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — printf on scalar",
			input:    `printf '%s\n' "$msg"`,
			expected: []katas.Violation{},
		},
		{
			name:  `invalid — printf '%s\n' "${array[@]}"`,
			input: `printf '%s\n' "${array[@]}"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1591",
					Message: "`printf '%s\\n' \"${array[@]}\"` — use Zsh `print -l -r -- \"${array[@]}\"` or `${(F)array}` for newline-joined output.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  `invalid — printf '%s' "${a[*]}"`,
			input: `printf '%s' "${a[*]}"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1591",
					Message: "`printf '%s\\n' \"${array[@]}\"` — use Zsh `print -l -r -- \"${array[@]}\"` or `${(F)array}` for newline-joined output.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1591")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
